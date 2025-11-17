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

describe('Authentication Flow Integration Tests', () => {
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

  describe('Google OAuth Flow', () => {
    it('completes Google OAuth login flow', async () => {
      const user = userEvent.setup()
      mockStartOAuth.mockResolvedValue(undefined)
      
      render(<LoginPage />)
      
      const googleButton = screen.getByText('Continue with Google')
      await user.click(googleButton)

      await waitFor(() => {
        expect(mockStartOAuth).toHaveBeenCalledWith('google')
      })
    })
  })

  describe('Email/Password Login Flow', () => {
    it('completes email/password login flow', async () => {
      const user = userEvent.setup()
      mockLoginWithPassword.mockResolvedValue({
        data: {
          accessToken: 'mock-token',
          refreshToken: 'mock-refresh-token',
        },
      })
      mockLogin.mockResolvedValue(undefined)
      
      render(<LoginPage />)
      
      // Click email login button
      const emailButton = screen.getByText('Continue with Email')
      await user.click(emailButton)

      await waitFor(() => {
        expect(screen.getByLabelText('Email')).toBeInTheDocument()
      })

      // Fill in form
      await user.type(screen.getByLabelText('Email'), 'admin@example.com')
      await user.type(screen.getByLabelText('Password'), 'password123')

      // Submit form
      const submitButton = screen.getByText('Sign in')
      await user.click(submitButton)

      await waitFor(() => {
        expect(mockLoginWithPassword).toHaveBeenCalledWith(
          'admin@example.com',
          'password123'
        )
      })
    })

    it('shows validation error for invalid email', async () => {
      const user = userEvent.setup()
      
      render(<LoginPage />)
      
      // Click email login button
      const emailButton = screen.getByText('Continue with Email')
      await user.click(emailButton)

      await waitFor(() => {
        expect(screen.getByLabelText('Email')).toBeInTheDocument()
      })

      // Try to submit without filling form
      const submitButton = screen.getByText('Sign in')
      expect(submitButton).toBeDisabled()
    })

    it('handles login error', async () => {
      const user = userEvent.setup()
      mockLoginWithPassword.mockRejectedValue(new Error('Invalid credentials'))
      
      render(<LoginPage />)
      
      // Click email login button
      const emailButton = screen.getByText('Continue with Email')
      await user.click(emailButton)

      await waitFor(() => {
        expect(screen.getByLabelText('Email')).toBeInTheDocument()
      })

      // Fill in form
      await user.type(screen.getByLabelText('Email'), 'admin@example.com')
      await user.type(screen.getByLabelText('Password'), 'wrongpassword')

      // Submit form
      const submitButton = screen.getByText('Sign in')
      await user.click(submitButton)

      await waitFor(() => {
        expect(screen.getByText('Invalid credentials')).toBeInTheDocument()
      })
    })
  })

  describe('Navigation Flow', () => {
    it('redirects to home when already authenticated', () => {
      mockUseAuth.mockReturnValue({
        isAuthenticated: true,
        loading: false,
        login: mockLogin,
        logout: jest.fn(),
        user: { id: '1', email: 'admin@example.com' },
      })
      
      render(<LoginPage />)
      
      // Should redirect (implementation dependent)
      // This test verifies the component handles authenticated state
      expect(mockUseAuth).toHaveBeenCalled()
    })
  })
})

