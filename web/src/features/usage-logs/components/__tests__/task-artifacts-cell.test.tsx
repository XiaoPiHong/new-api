/*
Copyright (C) 2023-2026 QuantumNous

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program. If not, see <https://www.gnu.org/licenses/>.

For commercial licensing, please contact support@quantumnous.com
*/
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import {
  act,
  cleanup,
  fireEvent,
  render,
  screen,
  waitFor,
} from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { afterEach, describe, expect, test, vi } from 'vitest'

import { api } from '@/lib/api'

import type { TaskLog } from '../../types'
import { TaskArtifactsCell } from '../task-artifacts'

let queryClient: QueryClient | undefined
afterEach(() => {
  cleanup()
  queryClient?.clear()
  vi.restoreAllMocks()
})

function contentUrl(expiresAt: number, key = 'image-main') {
  return `https://media.example/v1/tasks/task-public/artifacts/${key}/content?access=v2.${expiresAt}.${'A'.repeat(43)}`
}

function artifactResponse(url: string, legacy = false) {
  return {
    data: {
      success: true,
      data: legacy
        ? { artifacts: [], legacy_content_url: url }
        : {
            artifacts: [{ key: 'image-main', type: 'image', content_url: url }],
          },
    },
  }
}

function taskLog(overrides: Partial<TaskLog> = {}): TaskLog {
  return {
    id: 1,
    user_id: 7,
    platform: 'alibaba',
    task_id: 'task-public',
    action: 'text_to_image',
    channel_id: 357,
    group: 'default',
    quota: 25000,
    submit_time: 1,
    status: 'SUCCESS',
    admin_info: {
      task_plugin: {
        key: 'alibaba',
        name: 'Alibaba Bailian',
        version: '1.4.0',
      },
    },
    ...overrides,
  }
}

function renderCell(log: TaskLog) {
  queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false } },
  })
  return render(
    <QueryClientProvider client={queryClient}>
      <TaskArtifactsCell log={log} />
    </QueryClientProvider>
  )
}

describe('task artifacts cell', () => {
  test('media retry requests a new signed URL and disables repeated retries while metadata refreshes', async () => {
    const oldUrl = contentUrl(1900000900)
    const freshUrl = contentUrl(1900001800)
    let finishRefresh!: (response: ReturnType<typeof artifactResponse>) => void
    const get = vi
      .spyOn(api, 'get')
      .mockResolvedValueOnce(artifactResponse(oldUrl))
      .mockImplementationOnce(
        () =>
          new Promise((resolve) => {
            finishRefresh = resolve
          })
      )
    const user = userEvent.setup()
    renderCell(taskLog())
    await user.click(screen.getByRole('button', { name: 'Artifacts' }))
    const image = await screen.findByRole('img', { name: 'image-main' })
    expect(image).toHaveAttribute('src', oldUrl)
    fireEvent.error(image)
    const retry = screen.getByRole('button', { name: 'Retry' })
    await user.click(retry)
    await waitFor(() => expect(get).toHaveBeenCalledTimes(2))
    expect(retry).toBeDisabled()
    await act(async () => finishRefresh(artifactResponse(freshUrl)))
    expect(
      await screen.findByRole('img', { name: 'image-main' })
    ).toHaveAttribute('src', freshUrl)
  })

  test('a failed signed-URL refresh shows a retryable metadata error instead of replaying the expired media URL', async () => {
    const get = vi
      .spyOn(api, 'get')
      .mockResolvedValueOnce(artifactResponse(contentUrl(1900000900)))
      .mockRejectedValueOnce(new Error('Temporarily unavailable'))
    const user = userEvent.setup()
    renderCell(taskLog())
    await user.click(screen.getByRole('button', { name: 'Artifacts' }))
    fireEvent.error(await screen.findByRole('img', { name: 'image-main' }))
    await user.click(screen.getByRole('button', { name: 'Retry' }))
    expect(await screen.findByText('Failed to load artifacts')).toBeVisible()
    expect(
      screen.queryByRole('img', { name: 'image-main' })
    ).not.toBeInTheDocument()
    expect(get).toHaveBeenCalledTimes(2)
  })

  test('legacy video retry also obtains a replacement signed content URL', async () => {
    const oldUrl = contentUrl(1900000900, 'video')
    const freshUrl = contentUrl(1900001800, 'video')
    vi.spyOn(api, 'get')
      .mockResolvedValueOnce(artifactResponse(oldUrl, true))
      .mockResolvedValueOnce(artifactResponse(freshUrl, true))
    const user = userEvent.setup()
    const view = renderCell(
      taskLog({ admin_info: undefined, legacy_video_available: true })
    )
    await user.click(
      screen.getByRole('button', { name: 'Click to preview video' })
    )
    await waitFor(() =>
      expect(view.baseElement.querySelector('video')).toHaveAttribute(
        'src',
        oldUrl
      )
    )
    const video = view.baseElement.querySelector('video')
    if (!video) throw new Error('Video preview expected')
    fireEvent.error(video)
    await user.click(screen.getByRole('button', { name: 'Retry' }))
    await waitFor(() =>
      expect(view.baseElement.querySelector('video')).toHaveAttribute(
        'src',
        freshUrl
      )
    )
  })

  test('given a synchronous result the gateway did not retain, a hint replaces the artifacts button', () => {
    renderCell(taskLog({ result_discarded: true }))

    expect(screen.getByText('Result not retained')).toBeVisible()
    expect(
      screen.queryByRole('button', { name: /Artifacts/ })
    ).not.toBeInTheDocument()
  })

  test('given a retained successful plugin task, the artifacts button is offered', () => {
    renderCell(taskLog())

    expect(screen.getByRole('button', { name: /Artifacts/ })).toBeVisible()
    expect(screen.queryByText('Result not retained')).not.toBeInTheDocument()
  })

  test('given an unfinished task, neither the button nor the hint is shown', () => {
    renderCell(taskLog({ status: 'IN_PROGRESS', result_discarded: true }))

    expect(screen.getByText('-')).toBeVisible()
    expect(screen.queryByText('Result not retained')).not.toBeInTheDocument()
    expect(
      screen.queryByRole('button', { name: /Artifacts/ })
    ).not.toBeInTheDocument()
  })
})
