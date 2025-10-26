import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

class ForgotPasswordPage extends StatefulWidget {
  const ForgotPasswordPage({super.key});
  
  @override
  State<ForgotPasswordPage> createState() => _ForgotPasswordPageState();
}

class _ForgotPasswordPageState extends State<ForgotPasswordPage> {
  final _email = TextEditingController();

  @override
  void dispose() {
    _email.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: const Color(0xFFF5F5F5),
      body: SafeArea(
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            // Back button
            Padding(
              padding: const EdgeInsets.all(16),
              child: IconButton(
                icon: const Icon(
                  Icons.arrow_back_ios,
                  color: Colors.black,
                  size: 20,
                ),
                onPressed: () => context.pop(),
                padding: EdgeInsets.zero,
                constraints: const BoxConstraints(),
              ),
            ),
            
            Expanded(
              child: SingleChildScrollView(
                padding: const EdgeInsets.symmetric(horizontal: 24),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    const SizedBox(height: 20),
                    
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
                      'Forgot Password',
                      style: TextStyle(
                        fontSize: 26,
                        fontWeight: FontWeight.w600,
                        color: Colors.black,
                      ),
                    ),
                    
                    const SizedBox(height: 8),
                    
                    // Subtitle
                    const Text(
                      'Enter your email to reset password',
                      style: TextStyle(
                        fontSize: 15,
                        color: Colors.grey,
                        fontWeight: FontWeight.w400,
                      ),
                    ),
                    
                    const SizedBox(height: 50),
                    
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
                      decoration: const InputDecoration(
                        hintText: 'imshuvo97@gmail.com',
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
                    
                    const SizedBox(height: 50),
                    
                    // Reset Password button
                    SizedBox(
                      width: double.infinity,
                      height: 56,
                      child: ElevatedButton(
                        onPressed: () {
                          // Add reset password logic here
                          // Show success message or navigate
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
                          'Reset Password',
                          style: TextStyle(
                            fontSize: 16,
                            fontWeight: FontWeight.w600,
                          ),
                        ),
                      ),
                    ),
                    
                    const SizedBox(height: 60),
                  ],
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }
}