/*
Copyright (C) 2023-2026 QuantumNous

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program. If not, see <https://www.gnu.org/licenses/>.
*/
import { useEffect, useState } from 'react'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { ShieldCheck } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { toast } from 'sonner'
import { Input } from '@/components/ui/input'
import { SettingsSwitchField } from '../components/settings-form-layout'
import { SettingsPageFormActions } from '../components/settings-page-context'
import { SettingsForm } from '../components/settings-form-layout'
import { SettingsSection } from '../components/settings-section'
import { getRelayServiceAuth, updateRelayServiceAuth } from '../api'

export function RelayServiceAuthSection() {
  const { t } = useTranslation()
  const queryClient = useQueryClient()
  const configQuery = useQuery({
    queryKey: ['relay-service-auth'],
    queryFn: getRelayServiceAuth,
    staleTime: 0,
  })
  const updateMutation = useMutation({
    mutationFn: updateRelayServiceAuth,
    onSuccess: async (data) => {
      if (!data.success) {
        toast.error(data.message || t('Failed to update setting'))
        return
      }
      setSecret('')
      await queryClient.invalidateQueries({
        queryKey: ['relay-service-auth'],
      })
      toast.success(t('Setting updated successfully'))
    },
    onError: (error: Error) => {
      toast.error(error.message || t('Failed to update setting'))
    },
  })

  const [enabled, setEnabled] = useState(false)
  const [secret, setSecret] = useState('')
  const [savedEnabled, setSavedEnabled] = useState(false)
  const [secretConfigured, setSecretConfigured] = useState(false)

  const serverConfig = configQuery.data?.data
  useEffect(() => {
    if (!serverConfig) return
    setEnabled(serverConfig.enabled)
    setSavedEnabled(serverConfig.enabled)
    setSecretConfigured(serverConfig.secret_configured)
    setSecret('')
  }, [serverConfig])

  const isDirty = enabled !== savedEnabled || secret.trim() !== ''
  const isSaveDisabled =
    !isDirty || (enabled && !secretConfigured && secret.trim() === '')

  const handleSave = () => {
    const nextSecret = secret.trim()
    updateMutation.mutate({
      enabled,
      ...(nextSecret ? { secret: nextSecret } : {}),
    })
  }

  return (
    <SettingsSection title={t('Relay Service Authentication')}>
      {configQuery.isLoading ? (
        <div className='text-muted-foreground flex min-h-32 items-center justify-center text-sm'>
          {t('Loading settings...')}
        </div>
      ) : configQuery.isError ? (
        <div className='text-destructive text-sm'>
          {t('Failed to load settings')}
        </div>
      ) : (
        <SettingsForm
          autoComplete='off'
          onSubmit={(event) => {
            event.preventDefault()
            handleSave()
          }}
        >
          <SettingsPageFormActions
            onSave={handleSave}
            isSaving={updateMutation.isPending}
            isSaveDisabled={isSaveDisabled}
            saveLabel='Save relay settings'
          />

          <SettingsSwitchField
            checked={enabled}
            onCheckedChange={setEnabled}
            label={t('Enable relay service authentication')}
            description={t(
              'Require X-Xphai-Relay-Secret for protected /v1 relay requests.'
            )}
            disabled={updateMutation.isPending}
          />

          <div className='min-w-0 space-y-2 lg:col-span-2'>
            <label className='text-sm font-medium' htmlFor='relay-service-secret'>
              {t('Relay service secret')}
            </label>
            <Input
              id='relay-service-secret'
              type='password'
              value={secret}
              onChange={(event) => setSecret(event.target.value)}
              placeholder={
                secretConfigured
                  ? t('Leave blank to keep the existing secret')
                  : t('Enter a shared secret')
              }
              autoComplete='new-password'
              disabled={updateMutation.isPending}
            />
            <div className='text-muted-foreground flex items-start gap-2 text-xs'>
              <ShieldCheck className='mt-0.5 size-4 shrink-0' />
              <span>
                {secretConfigured
                  ? t('A relay service secret is configured.')
                  : t('No relay service secret is configured.')}{' '}
                {t(
                  'The secret is never returned to the browser. Keep the same value in xphai-web and SSO.'
                )}
              </span>
            </div>
          </div>

          <p className='text-muted-foreground text-xs lg:col-span-2'>
            {t(
              'Changing the secret takes effect immediately. Update xphai-web and SSO at the same time or their requests will receive 403.'
            )}
          </p>
        </SettingsForm>
      )}
    </SettingsSection>
  )
}
