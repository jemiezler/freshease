import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

class SignupPage extends StatefulWidget {
  const SignupPage({super.key});
  
  @override
  State<SignupPage> createState() => _SignupPageState();
}

class _SignupPageState extends State<SignupPage> {
  final _username = TextEditingController();
  final _email = TextEditingController();
  final _password = TextEditingController();
  bool _obscurePassword = true;
  bool _isEmailValid = false;

  @override
  void initState() {
    super.initState();
    _email.addListener(_validateEmail);
  }

  void _validateEmail() {
    final email = _email.text;
    final isValid = RegExp(r'^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$').hasMatch(email);
    setState(() {
      _isEmailValid = isValid && email.isNotEmpty;
    });
  }

  @override
  void dispose() {
    _username.dispose();
    _email.dispose();
    _password.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: const Color(0xFFF5F5F5),
      body: SafeArea(
        child: SingleChildScrollView(
          padding: const EdgeInsets.symmetric(horizontal: 24),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              const SizedBox(height: 60),
              
              // Logo
              Center(
                child: Image.asset(
                  'lib/assets/logo.png',
                  width: 70,
                  height: 70,
                ),
              ),
              
              const SizedBox(height: 80),
              
              // Title
              const Text(
                'Sign Up',
                style: TextStyle(
                  fontSize: 26,
                  fontWeight: FontWeight.w600,
                  color: Colors.black,
                ),
              ),
              
              const SizedBox(height: 8),
              
              // Subtitle
              const Text(
                'Enter your credentials to continue',
                style: TextStyle(
                  fontSize: 15,
                  color: Colors.grey,
                  fontWeight: FontWeight.w400,
                ),
              ),
              
              const SizedBox(height: 50),
              
              // Username label
              const Text(
                'Username',
                style: TextStyle(
                  fontSize: 15,
                  color: Colors.grey,
                  fontWeight: FontWeight.w400,
                ),
              ),
              
              const SizedBox(height: 12),
              
              // Username field
              TextField(
                controller: _username,
                decoration: const InputDecoration(
                  hintText: 'Afsar Hossen Shuvo',
                  hintStyle: TextStyle(
                    color: Colors.black,
                    fontSize: 16,
                  ),
                  border: UnderlineInputBorder(
                    borderSide: BorderSide(color: Color(0xFFE0E0E0)),
                  ),
                  enabledBorder: UnderlineInputBorder(
                    borderSide: BorderSide(color: Color(0xFFE0E0E0)),
                  ),
                  focusedBorder: UnderlineInputBorder(
                    borderSide: BorderSide(color: Color(0xFF5FBE7E), width: 2),
                  ),
                  contentPadding: EdgeInsets.symmetric(vertical: 12),
                ),
                style: const TextStyle(
                  fontSize: 16,
                  color: Colors.black,
                ),
              ),
              
              const SizedBox(height: 40),
              
              // Email label
              const Text(
                'Email',
                style: TextStyle(
                  fontSize: 15,
                  color: Colors.grey,
                  fontWeight: FontWeight.w400,
                ),
              ),
              
              const SizedBox(height: 12),
              
              // Email field
              TextField(
                controller: _email,
                keyboardType: TextInputType.emailAddress,
                decoration: InputDecoration(
                  hintText: 'imshuvo97@gmail.com',
                  hintStyle: const TextStyle(
                    color: Colors.black,
                    fontSize: 16,
                  ),
                  border: const UnderlineInputBorder(
                    borderSide: BorderSide(color: Color(0xFFE0E0E0)),
                  ),
                  enabledBorder: const UnderlineInputBorder(
                    borderSide: BorderSide(color: Color(0xFFE0E0E0)),
                  ),
                  focusedBorder: const UnderlineInputBorder(
                    borderSide: BorderSide(color: Color(0xFF5FBE7E), width: 2),
                  ),
                  contentPadding: const EdgeInsets.symmetric(vertical: 12),
                  suffixIcon: _isEmailValid
                      ? const Icon(
                          Icons.check,
                          color: Color(0xFF5FBE7E),
                          size: 22,
                        )
                      : null,
                ),
                style: const TextStyle(
                  fontSize: 16,
                  color: Colors.black,
                ),
              ),
              
              const SizedBox(height: 40),
              
              // Password label
              const Text(
                'Password',
                style: TextStyle(
                  fontSize: 15,
                  color: Colors.grey,
                  fontWeight: FontWeight.w400,
                ),
              ),
              
              const SizedBox(height: 12),
              
              // Password field
              TextField(
                controller: _password,
                obscureText: _obscurePassword,
                decoration: InputDecoration(
                  hintText: '••••••••',
                  hintStyle: const TextStyle(
                    color: Colors.black,
                    fontSize: 16,
                  ),
                  border: const UnderlineInputBorder(
                    borderSide: BorderSide(color: Color(0xFFE0E0E0)),
                  ),
                  enabledBorder: const UnderlineInputBorder(
                    borderSide: BorderSide(color: Color(0xFFE0E0E0)),
                  ),
                  focusedBorder: const UnderlineInputBorder(
                    borderSide: BorderSide(color: Color(0xFF5FBE7E), width: 2),
                  ),
                  contentPadding: const EdgeInsets.symmetric(vertical: 12),
                  suffixIcon: IconButton(
                    icon: Icon(
                      _obscurePassword ? Icons.visibility_off_outlined : Icons.visibility_outlined,
                      color: Colors.grey,
                      size: 22,
                    ),
                    onPressed: () {
                      setState(() {
                        _obscurePassword = !_obscurePassword;
                      });
                    },
                  ),
                ),
                style: const TextStyle(
                  fontSize: 16,
                  color: Colors.black,
                ),
              ),
              
              const SizedBox(height: 24),
              
              // Terms and privacy text
              RichText(
                text: const TextSpan(
                  style: TextStyle(
                    fontSize: 13,
                    color: Colors.grey,
                    height: 1.5,
                  ),
                  children: [
                    TextSpan(text: 'By continuing you agree to our '),
                    TextSpan(
                      text: 'Terms of Service',
                      style: TextStyle(
                        color: Color(0xFF5FBE7E),
                      ),
                    ),
                    TextSpan(text: ' and '),
                    TextSpan(
                      text: 'Privacy Policy',
                      style: TextStyle(
                        color: Color(0xFF5FBE7E),
                      ),
                    ),
                    TextSpan(text: '.'),
                  ],
                ),
              ),
              
              const SizedBox(height: 32),
              
              // Sign Up button
              SizedBox(
                width: double.infinity,
                height: 56,
                child: ElevatedButton(
                  onPressed: () {
                    // Add signup logic here
                    context.go('/');
                  },
                  style: ElevatedButton.styleFrom(
                    backgroundColor: const Color(0xFF5FBE7E),
                    foregroundColor: Colors.white,
                    shape: RoundedRectangleBorder(
                      borderRadius: BorderRadius.circular(12),
                    ),
                    elevation: 0,
                  ),
                  child: const Text(
                    'Sign Up',
                    style: TextStyle(
                      fontSize: 16,
                      fontWeight: FontWeight.w600,
                    ),
                  ),
                ),
              ),
              
              const SizedBox(height: 24),
              
              // Login text
              Center(
                child: Row(
                  mainAxisSize: MainAxisSize.min,
                  children: [
                    const Text(
                      'Already have an account? ',
                      style: TextStyle(
                        color: Colors.black,
                        fontSize: 14,
                        fontWeight: FontWeight.w400,
                      ),
                    ),
                    GestureDetector(
                      onTap: () {
                        context.go('/login');
                      },
                      child: const Text(
                        'Signin',
                        style: TextStyle(
                          color: Color(0xFF5FBE7E),
                          fontSize: 14,
                          fontWeight: FontWeight.w500,
                        ),
                      ),
                    ),
                  ],
                ),
              ),
              
              const SizedBox(height: 60),
            ],
          ),
        ),
      ),
    );
  }
}