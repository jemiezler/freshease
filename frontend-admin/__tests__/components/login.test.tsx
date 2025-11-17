import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import LoginPage from '@/app/login/page'
import { useAuth } from '@/lib/auth-context'
import { startOAuth, loginWithPassword } from '@/lib/auth'

// Mock dependencies
jest.mock('@/lib/auth-context')
jest.mock('@/lib/auth')
jest.mock('next/navigation', () => ({
  useRouter: () => ({
    push: jest.fn(),
  }),
  useSearchParams: () => new URLSearchParams(),
}))

const mockUseAuth = useAuth as jest.MockedFunction<typeof useAuth>
const mockStartOAuth = startOAuth as jest.MockedFunction<typeof startOAuth>
const mockLoginWithPassword = loginWithPassword as jest.MockedFunction<typeof loginWithPassword>

describe('LoginPage', () => {
  const mockLogin = jest.fn()
  const mockPush = jest.fn()

  beforeEach(() => {
    // Clear individual mocks
    mockLogin.mockClear()
    mockPush.mockClear()
    mockUseAuth.mockClear()
    mockStartOAuth.mockClear()
    mockLoginWithPassword.mockClear()
    
    mockUseAuth.mockReturnValue({
      isAuthenticated: false,
      loading: false,
      login: mockLogin,
      logout: jest.fn(),
      user: null,
    })
  })

  it('renders login page with title', () => {
    render(<LoginPage />)
    expect(screen.getByText('Freshease Admin')).toBeInTheDocument()
    expect(screen.getByText('Sign in to access the admin panel')).toBeInTheDocument()
  })

  it('displays Google login button', () => {
    render(<LoginPage />)
    expect(screen.getByText('Continue with Google')).toBeInTheDocument()
  })

  it('displays email login option', () => {
    render(<LoginPage />)
    expect(screen.getByText('Continue with Email')).toBeInTheDocument()
  })

  it('shows loading state when loading', () => {
    mockUseAuth.mockReturnValue({
      isAuthenticated: false,
      loading: true,
      login: mockLogin,
      logout: jest.fn(),
      user: null,
    })
    render(<LoginPage />)
    expect(screen.getByText('Loading...')).toBeInTheDocument()
  })

  it('shows email/password form when email button is clicked', async () => {
    const user = userEvent.setup()
    render(<LoginPage />)
    
    const emailButton = screen.getByText('Continue with Email')
    await user.click(emailButton)

    await waitFor(() => {
      expect(screen.getByLabelText('Email')).toBeInTheDocument()
      expect(screen.getByLabelText('Password')).toBeInTheDocument()
    })
  })

  it('calls startOAuth when Google login button is clicked', async () => {
    const user = userEvent.setup()
    mockStartOAuth.mockResolvedValue(undefined)
    
    render(<LoginPage />)
    
    const googleButton = screen.getByText('Continue with Google')
    await user.click(googleButton)

    await waitFor(() => {
      expect(mockStartOAuth).toHaveBeenCalledWith('google')
    })
  })

  it('shows error message when authentication fails', () => {
    // Mock URLSearchParams to return error
    jest.spyOn(URLSearchParams.prototype, 'get').mockReturnValue('auth_failed')
    
    render(<LoginPage />)
    
    expect(screen.getByText(/Authentication failed/i)).toBeInTheDocument()
  })

  it('displays init admin button', () => {
    render(<LoginPage />)
    expect(screen.getByText('Initialize Admin User')).toBeInTheDocument()
  })
})

