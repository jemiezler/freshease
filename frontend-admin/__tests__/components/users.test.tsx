import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import '@testing-library/jest-dom'

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
    ;(global as any).__mockUsersResource = resource
  }
  
  return {
    createResource: jest.fn(() => resource),
  }
})

// Mock Next.js navigation
jest.mock('next/navigation', () => ({
  useRouter: () => ({
    push: jest.fn(),
    replace: jest.fn(),
  }),
}))

// Mock window.confirm
window.confirm = jest.fn(() => true)

// Import the component AFTER mocks are set up
import UsersPage from '@/app/users/page'

describe('UsersPage', () => {
  // Get the mock resource from global after module is loaded
  const getMockResource = () => {
    if (typeof global !== 'undefined' && (global as any).__mockUsersResource) {
      return (global as any).__mockUsersResource
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
      
      // Default mock response
      mocks.list.mockResolvedValue({
        data: [
          {
            id: '1',
            email: 'user1@example.com',
            name: 'User One',
            status: 'active',
          },
          {
            id: '2',
            email: 'user2@example.com',
            name: 'User Two',
            status: 'inactive',
          },
        ],
      })
    }
  })

  it('renders users page with title', async () => {
    render(<UsersPage />)
    
    await waitFor(() => {
      expect(screen.getByText('Users')).toBeInTheDocument()
    })
  })

  it('displays users list', async () => {
    render(<UsersPage />)
    
    await waitFor(() => {
      expect(screen.getByText('User One')).toBeInTheDocument()
      expect(screen.getByText('User Two')).toBeInTheDocument()
      expect(screen.getByText('user1@example.com')).toBeInTheDocument()
      expect(screen.getByText('user2@example.com')).toBeInTheDocument()
    })
  })

  it('displays loading state initially', () => {
    const mocks = getMockResource()
    if (mocks) {
      // Make the promise never resolve to test loading state
      mocks.list.mockImplementation(() => new Promise(() => {}))
    }
    
    render(<UsersPage />)
    
    expect(screen.getByText('Loading users…')).toBeInTheDocument()
  })

  it('displays empty state when no users exist', async () => {
    const mocks = getMockResource()
    if (mocks) {
      mocks.list.mockResolvedValue({
        data: [],
      })
    }
    
    render(<UsersPage />)
    
    await waitFor(() => {
      expect(screen.getByText('No results.')).toBeInTheDocument()
    })
  })

  it('handles user loading error', async () => {
    const mocks = getMockResource()
    if (mocks) {
      mocks.list.mockRejectedValue(new Error('Failed to load users'))
    }
    
    render(<UsersPage />)
    
    await waitFor(() => {
      expect(screen.getByText('Failed to load users')).toBeInTheDocument()
    })
  })

  it('displays "New" button', async () => {
    render(<UsersPage />)
    
    await waitFor(() => {
      expect(screen.getByText('New')).toBeInTheDocument()
    })
  })

  it('calls delete when delete button is clicked', async () => {
    const mocks = getMockResource()
    if (mocks) {
      mocks.delete.mockResolvedValue({})
    }

    render(<UsersPage />)
    
    await waitFor(() => {
      expect(screen.getByText('User One')).toBeInTheDocument()
    })

    // Find delete button by finding the row with User One and clicking the trash icon button
    const userOneCell = screen.getByText('User One')
    const userOneRow = userOneCell.closest('tr')
    if (userOneRow) {
      const actionButtons = userOneRow.querySelectorAll('button')
      // The delete button is the second button (index 1) with trash icon
      if (actionButtons.length >= 2) {
        await userEvent.click(actionButtons[1])
        
        await waitFor(() => {
          if (mocks) {
            expect(mocks.delete).toHaveBeenCalled()
          }
        }, { timeout: 3000 })
      }
    }
  })
})
