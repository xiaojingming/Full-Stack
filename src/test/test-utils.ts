import { render as rtlRender, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { type ReactElement } from 'react'

function render(ui: ReactElement) {
  return {
    user: userEvent.setup(),
    ...rtlRender(ui),
  }
}

export { render, screen, waitFor, userEvent }
