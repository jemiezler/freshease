import { ReactElement } from 'react'
import { render, RenderOptions } from '@testing-library/react'
import { AuthProvider } from '@/lib/auth-context'

// Custom render function that includes providers
const AllTheProviders = ({ children }: { children: React.ReactNode }) => {
  return (
    <AuthProvider>
      {children}
    </AuthProvider>
  )
}

const customRender = (
  ui: ReactElement,
  options?: Omit<RenderOptions, 'wrapper'>,
) => render(ui, { wrapper: AllTheProviders, ...options })

export * from '@testing-library/react'
export { customRender as render }

// Mock data helpers
export const mockUser = {
  id: '1',
  email: 'admin@example.com',
  name: 'Admin User',
}

export const mockProduct = {
  id: '1',
  name: 'Test Product',
  price: 99.99,
  sku: 'TEST-001',
  description: 'Test description',
  unitLabel: 'kg',
  isActive: true,
  imageUrl: 'https://example.com/image.jpg',
  createdAt: new Date().toISOString(),
  updatedAt: new Date().toISOString(),
}

export const mockCategory = {
  id: '1',
  name: 'Test Category',
  slug: 'test-category',
  description: 'Test category description',
  createdAt: new Date().toISOString(),
  updatedAt: new Date().toISOString(),
}

