import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'

// Create mock functions that will be shared between the factory and tests
// Use a module-level object that gets populated in the factory
let mockResource: {
  list: jest.Mock
  create: jest.Mock
  update: jest.Mock
  delete: jest.Mock
  get: jest.Mock
}

// Mock the resource module
jest.mock('@/lib/resource', () => {
  const list = jest.fn()
  const create = jest.fn()
  const update = jest.fn()
  const del = jest.fn()
  const get = jest.fn()
  
  const resource = {
    list,
    create,
    update,
    delete: del,
    get,
  }
  
  // Store in module-level variable for test access
  // This assignment happens when the mock is created
  if (typeof global !== 'undefined') {
    (global as any).__mockResource = resource
  }
  
  return {
    createResource: jest.fn(() => resource),
  }
})

// Mock Next.js navigation
jest.mock('next/navigation', () => ({
  useRouter: () => ({
    push: jest.fn(),
  }),
}))

// Mock window.confirm
window.confirm = jest.fn(() => true)

// Import the component AFTER mocks are set up
import ProductsPage from '@/app/products/page'

describe('ProductsPage', () => {
  // Get the mock resource from global after module is loaded
  const getMockResource = () => {
    if (typeof global !== 'undefined' && (global as any).__mockResource) {
      return (global as any).__mockResource
    }
    // Fallback: access through the mocked module
    const resourceModule = require('@/lib/resource')
    const createResource = resourceModule.createResource as jest.Mock
    return createResource.mock.results[0].value
  }

  beforeEach(() => {
    const mocks = getMockResource()
    if (mocks) {
      // Reset all mocks before each test
      mocks.list.mockClear()
      mocks.delete.mockClear()
      mocks.create.mockClear()
      mocks.update.mockClear()
      mocks.get.mockClear()
      
      // Default mock responses - categories list returns empty
      mocks.list.mockResolvedValue({
        data: [],
      })
    }
  })

  it('renders products page with title', async () => {
    render(<ProductsPage />)
    
    await waitFor(() => {
      expect(screen.getByText('Products')).toBeInTheDocument()
    })
  })

  it('displays loading state initially', () => {
    const mocks = getMockResource()
    if (mocks) {
      // Make the promise never resolve to test loading state
      mocks.list.mockImplementation(() => new Promise(() => {}))
    }
    
    render(<ProductsPage />)
    
    expect(screen.getByText('Loading…')).toBeInTheDocument()
  })

  it('displays "No categories found" when no categories exist', async () => {
    const mocks = getMockResource()
    if (mocks) {
      mocks.list.mockResolvedValue({
        data: [],
      })
    }
    
    render(<ProductsPage />)
    
    await waitFor(() => {
      expect(screen.getByText('No categories found.')).toBeInTheDocument()
    })
  })

  it('displays "New Category" button', async () => {
    render(<ProductsPage />)
    
    await waitFor(() => {
      expect(screen.getByText('New Category')).toBeInTheDocument()
    })
  })

  it('handles category loading error', async () => {
    const mocks = getMockResource()
    if (mocks) {
      mocks.list.mockRejectedValue(new Error('Failed to load'))
    }
    
    render(<ProductsPage />)
    
    await waitFor(() => {
      expect(screen.getByText('Failed to load')).toBeInTheDocument()
    })
  })
})
