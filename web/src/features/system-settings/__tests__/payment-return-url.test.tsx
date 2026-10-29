import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { cleanup, render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import type { ComponentProps } from 'react'
import { afterEach, expect, test, vi } from 'vitest'

import { api } from '@/lib/api'

import { SettingsPageProvider } from '../components/settings-page-context'
import { PaymentSettingsSection } from '../integrations/payment-settings-section'

const props: ComponentProps<typeof PaymentSettingsSection> = {
  defaultValues: {
    TopUpReturnURL: 'https://ai.example.com/topup?payment=return',
    PayAddress: '',
    EpayId: '',
    EpayKey: '',
    Price: 7.3,
    MinTopUp: 1,
    CustomCallbackAddress: '',
    PayMethods: '[]',
    AmountOptions: '[]',
    AmountDiscount: '{}',
    StripeApiSecret: '',
    StripeWebhookSecret: '',
    StripePriceId: '',
    StripeUnitPrice: 1,
    StripeMinTopUp: 1,
    StripePromotionCodesEnabled: false,
    CreemApiKey: '',
    CreemWebhookSecret: '',
    CreemTestMode: false,
    CreemProducts: '[]',
  },
  waffoDefaultValues: {
    WaffoEnabled: false,
    WaffoApiKey: '',
    WaffoPrivateKey: '',
    WaffoPublicCert: '',
    WaffoSandboxPublicCert: '',
    WaffoSandboxApiKey: '',
    WaffoSandboxPrivateKey: '',
    WaffoSandbox: false,
    WaffoMerchantId: '',
    WaffoCurrency: 'USD',
    WaffoUnitPrice: 1,
    WaffoMinTopUp: 1,
    WaffoNotifyUrl: '',
    WaffoReturnUrl: '',
    WaffoPayMethods: '[]',
  },
  waffoPancakeDefaultValues: {
    WaffoPancakeMerchantID: '',
    WaffoPancakePrivateKey: '',
    WaffoPancakeReturnURL: '',
  },
  complianceDefaults: {
    confirmed: true,
    termsVersion: 'v1',
    confirmedAt: 1,
    confirmedBy: 1,
  },
}

afterEach(() => {
  cleanup()
  vi.restoreAllMocks()
})

function renderPaymentSettings() {
  const actions = document.createElement('div')
  document.body.append(actions)
  const client = new QueryClient({
    defaultOptions: { mutations: { retry: false } },
  })
  render(
    <QueryClientProvider client={client}>
      <SettingsPageProvider actionsContainer={actions}>
        <PaymentSettingsSection {...props} />
      </SettingsPageProvider>
    </QueryClientProvider>
  )
  return () => {
    client.clear()
    actions.remove()
  }
}

test('loads the payment return URL and persists a trimmed replacement using the namespaced option', async () => {
  const put = vi
    .spyOn(api, 'put')
    .mockResolvedValue({ data: { success: true } })
  const dispose = renderPaymentSettings()
  try {
    const user = userEvent.setup()
    const field = screen.getByRole('textbox', { name: 'Payment return URL' })
    expect(field).toHaveValue(props.defaultValues.TopUpReturnURL)
    await user.clear(field)
    await user.type(field, '  https://ai.example.com/account?payment=done  ')
    await user.click(screen.getByRole('button', { name: 'Save all settings' }))
    await waitFor(() =>
      expect(put).toHaveBeenCalledExactlyOnceWith('/api/option/', {
        key: 'payment_setting.top_up_return_url',
        value: 'https://ai.example.com/account?payment=done',
      })
    )
  } finally {
    dispose()
  }
})

test('blocks a non-HTTP payment return URL before updating any system options', async () => {
  const put = vi
    .spyOn(api, 'put')
    .mockResolvedValue({ data: { success: true } })
  const dispose = renderPaymentSettings()
  try {
    const user = userEvent.setup()
    const field = screen.getByRole('textbox', { name: 'Payment return URL' })
    await user.clear(field)
    await user.type(field, 'javascript:alert(1)')
    await user.click(screen.getByRole('button', { name: 'Save all settings' }))
    expect(
      await screen.findByText(
        'Provide a valid URL starting with http:// or https://'
      )
    ).toBeVisible()
    expect(put).not.toHaveBeenCalled()
  } finally {
    dispose()
  }
})
