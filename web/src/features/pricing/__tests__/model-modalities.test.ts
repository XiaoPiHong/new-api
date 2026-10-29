import { expect, it } from 'vitest'

import { getModelModalities } from '../lib/model-helpers'
import type { PricingModel } from '../types'

function model(overrides: Partial<PricingModel>): PricingModel {
  return {
    id: 1,
    model_name: 'image-model',
    quota_type: 1,
    model_ratio: 1,
    completion_ratio: 1,
    enable_groups: [],
    ...overrides,
  }
}

it('shows image editing input and output modalities from its declared endpoint', () => {
  expect(
    getModelModalities(model({ supported_endpoint_types: ['image-edit'] }))
  ).toEqual({ input: ['text', 'image'], output: ['image'] })
})

it('shows text to image for generation without assuming an image input', () => {
  expect(
    getModelModalities(
      model({ supported_endpoint_types: ['image-generation'] })
    )
  ).toEqual({ input: ['text'], output: ['image'] })
})

it('keeps explicit catalog modalities including an explicitly empty list', () => {
  expect(
    getModelModalities(
      model({
        supported_endpoint_types: ['image-edit'],
        input_modalities: [],
        output_modalities: ['file'],
      })
    )
  ).toEqual({ input: [], output: ['file'] })
})

it('keeps text output when the model also declares a chat endpoint', () => {
  expect(
    getModelModalities(
      model({ supported_endpoint_types: ['openai', 'image-edit'] })
    )
  ).toEqual({ input: ['text', 'image'], output: ['text', 'image'] })
})

it('does not infer modalities from an image-like model name without an image endpoint', () => {
  expect(
    getModelModalities(model({ model_name: 'image-vision-model' }))
  ).toEqual({
    input: [],
    output: [],
  })
})
