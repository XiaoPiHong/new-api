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
import { afterEach, expect, test, vi } from 'vitest'

import { api } from '@/lib/api'

import { getRelayServiceAuth, updateRelayServiceAuth } from '../api'

afterEach(() => vi.restoreAllMocks())

test('rejects a denied configuration read instead of exposing an empty editable form', async () => {
  vi.spyOn(api, 'get').mockResolvedValue({
    data: { success: false, message: 'Access denied' },
  })

  await expect(getRelayServiceAuth()).rejects.toThrow('Access denied')
})

test('rejects a denied secret update instead of reporting a saved setting', async () => {
  vi.spyOn(api, 'put').mockResolvedValue({
    data: { success: false, message: 'Access denied' },
  })

  await expect(
    updateRelayServiceAuth({ enabled: true, secret: 'test-only-secret' })
  ).rejects.toThrow('Access denied')
})

test('can toggle authentication without replacing an already configured secret', async () => {
  const response = {
    success: true,
    message: '',
    data: { enabled: false, secret_configured: true },
  }
  const put = vi.spyOn(api, 'put').mockResolvedValue({ data: response })

  await expect(updateRelayServiceAuth({ enabled: false })).resolves.toEqual(
    response
  )
  expect(put).toHaveBeenCalledWith('/api/relay-service-auth', {
    enabled: false,
  })
})
