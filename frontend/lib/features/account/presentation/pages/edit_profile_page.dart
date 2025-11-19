import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:go_router/go_router.dart';
import '../../domain/entities/user_profile.dart';
import '../state/user_cubit.dart';

class EditProfilePage extends StatefulWidget {
  final UserProfile user;

  const EditProfilePage({super.key, required this.user});

  @override
  State<EditProfilePage> createState() => _EditProfilePageState();
}

class _EditProfilePageState extends State<EditProfilePage> {
  late final TextEditingController _nameController;
  late final TextEditingController _phoneController;
  late final TextEditingController _bioController;
  late final TextEditingController _avatarController;
  late final TextEditingController _heightController;
  late final TextEditingController _weightController;
  DateTime? _selectedDate;
  String? _selectedSex;
  String? _selectedGoal;
  final _formKey = GlobalKey<FormState>();

  @override
  void initState() {
    super.initState();
    _nameController = TextEditingController(text: widget.user.name);
    _phoneController = TextEditingController(text: widget.user.phone ?? '');
    _bioController = TextEditingController(text: widget.user.bio ?? '');
    _avatarController = TextEditingController(text: widget.user.avatar ?? '');
    _heightController = TextEditingController(
      text: widget.user.heightCm?.toString() ?? '',
    );
    _weightController = TextEditingController(
      text: widget.user.weightKg?.toString() ?? '',
    );
    _selectedDate = widget.user.dateOfBirth;
    _selectedSex = widget.user.sex;
    _selectedGoal = widget.user.goal;
  }

  @override
  void dispose() {
    _nameController.dispose();
    _phoneController.dispose();
    _bioController.dispose();
    _avatarController.dispose();
    _heightController.dispose();
    _weightController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return BlocListener<UserCubit, UserState>(
      listener: (context, state) {
        if (state.isUpdating == false && state.error == null) {
          // Update successful, navigate back
          context.pop();
        } else if (state.error != null) {
          ScaffoldMessenger.of(context).showSnackBar(
            SnackBar(content: Text(state.error!), backgroundColor: Colors.red),
          );
        }
      },
      child: Scaffold(
        appBar: AppBar(
          title: const Text('Edit Profile'),
          actions: [
            BlocBuilder<UserCubit, UserState>(
              builder: (context, state) {
                return TextButton(
                  onPressed: state.isUpdating ? null : _saveProfile,
                  child: state.isUpdating
                      ? const SizedBox(
                          width: 16,
                          height: 16,
                          child: CircularProgressIndicator(strokeWidth: 2),
                        )
                      : const Text('Save'),
                );
              },
            ),
          ],
        ),
        body: SingleChildScrollView(
          padding: const EdgeInsets.all(16),
          child: Form(
            key: _formKey,
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                // Profile Picture Section
                Center(
                  child: Column(
                    children: [
                      CircleAvatar(
                        radius: 60,
                        backgroundImage: widget.user.avatar != null
                            ? NetworkImage(widget.user.avatar!)
                            : null,
                        child: widget.user.avatar == null
                            ? Text(
                                widget.user.initials,
                                style: const TextStyle(
                                  fontSize: 32,
                                  fontWeight: FontWeight.w800,
                                ),
                              )
                            : null,
                      ),
                      const SizedBox(height: 16),
                      Text(
                        'Profile Picture',
                        style: Theme.of(context).textTheme.titleMedium,
                      ),
                      const SizedBox(height: 8),
                      Text(
                        'Add a profile picture URL below',
                        style: Theme.of(context).textTheme.bodySmall?.copyWith(
                          color: Colors.grey[600],
                        ),
                      ),
                    ],
                  ),
                ),
                const SizedBox(height: 32),

                // Form Fields
                TextFormField(
                  controller: _nameController,
                  decoration: const InputDecoration(
                    labelText: 'Name',
                    border: OutlineInputBorder(),
                    prefixIcon: Icon(Icons.person),
                  ),
                  validator: (value) {
                    if (value == null || value.trim().isEmpty) {
                      return 'Name is required';
                    }
                    if (value.trim().length < 2) {
                      return 'Name must be at least 2 characters';
                    }
                    if (value.trim().length > 100) {
                      return 'Name must be less than 100 characters';
                    }
                    return null;
                  },
                ),
                const SizedBox(height: 16),

                TextFormField(
                  controller: _phoneController,
                  decoration: const InputDecoration(
                    labelText: 'Phone',
                    border: OutlineInputBorder(),
                    prefixIcon: Icon(Icons.phone),
                  ),
                  keyboardType: TextInputType.phone,
                  validator: (value) {
                    if (value != null && value.trim().isNotEmpty) {
                      if (value.trim().length < 10) {
                        return 'Phone must be at least 10 characters';
                      }
                      if (value.trim().length > 20) {
                        return 'Phone must be less than 20 characters';
                      }
                    }
                    return null;
                  },
                ),
                const SizedBox(height: 16),

                TextFormField(
                  controller: _bioController,
                  decoration: const InputDecoration(
                    labelText: 'Bio',
                    border: OutlineInputBorder(),
                    prefixIcon: Icon(Icons.info),
                  ),
                  maxLines: 3,
                  validator: (value) {
                    if (value != null && value.trim().isNotEmpty) {
                      if (value.trim().length < 10) {
                        return 'Bio must be at least 10 characters';
                      }
                      if (value.trim().length > 500) {
                        return 'Bio must be less than 500 characters';
                      }
                    }
                    return null;
                  },
                ),
                const SizedBox(height: 16),

                TextFormField(
                  controller: _avatarController,
                  decoration: const InputDecoration(
                    labelText: 'Avatar URL',
                    border: OutlineInputBorder(),
                    prefixIcon: Icon(Icons.image),
                  ),
                  validator: (value) {
                    if (value != null && value.trim().isNotEmpty) {
                      if (value.trim().length < 10) {
                        return 'Avatar URL must be at least 10 characters';
                      }
                      if (value.trim().length > 200) {
                        return 'Avatar URL must be less than 200 characters';
                      }
                      final uri = Uri.tryParse(value.trim());
                      if (uri == null || !uri.hasAbsolutePath) {
                        return 'Please enter a valid URL';
                      }
                    }
                    return null;
                  },
                ),
                const SizedBox(height: 16),

                // Date of Birth Field
                InkWell(
                  onTap: _selectDate,
                  child: InputDecorator(
                    decoration: const InputDecoration(
                      labelText: 'Date of Birth',
                      border: OutlineInputBorder(),
                      prefixIcon: Icon(Icons.calendar_today),
                    ),
                    child: Text(
                      _selectedDate != null
                          ? '${_selectedDate!.day}/${_selectedDate!.month}/${_selectedDate!.year}'
                          : 'Select date',
                      style: TextStyle(
                        color: _selectedDate != null
                            ? Theme.of(context).textTheme.bodyLarge?.color
                            : Theme.of(context).hintColor,
                      ),
                    ),
                  ),
                ),
                const SizedBox(height: 16),

                // Sex Field
                DropdownButtonFormField<String>(
                  initialValue: _selectedSex,
                  decoration: const InputDecoration(
                    labelText: 'Sex',
                    border: OutlineInputBorder(),
                    prefixIcon: Icon(Icons.person_outline),
                  ),
                  items: const [
                    DropdownMenuItem(value: null, child: Text('Not specified')),
                    DropdownMenuItem(value: 'male', child: Text('Male')),
                    DropdownMenuItem(value: 'female', child: Text('Female')),
                    DropdownMenuItem(value: 'other', child: Text('Other')),
                  ],
                  onChanged: (value) {
                    setState(() {
                      _selectedSex = value;
                    });
                  },
                ),
                const SizedBox(height: 16),

                // Goal Field
                DropdownButtonFormField<String>(
                  value: _selectedGoal,
                  decoration: const InputDecoration(
                    labelText: 'Goal',
                    border: OutlineInputBorder(),
                    prefixIcon: Icon(Icons.flag),
                  ),
                  items: const [
                    DropdownMenuItem(value: null, child: Text('Not specified')),
                    DropdownMenuItem(
                      value: 'maintenance',
                      child: Text('Maintenance'),
                    ),
                    DropdownMenuItem(
                      value: 'weight_loss',
                      child: Text('Weight Loss'),
                    ),
                    DropdownMenuItem(
                      value: 'weight_gain',
                      child: Text('Weight Gain'),
                    ),
                  ],
                  onChanged: (value) {
                    setState(() {
                      _selectedGoal = value;
                    });
                  },
                ),
                const SizedBox(height: 16),

                // Height and Weight Fields
                Row(
                  children: [
                    Expanded(
                      child: TextFormField(
                        controller: _heightController,
                        decoration: const InputDecoration(
                          labelText: 'Height (cm)',
                          border: OutlineInputBorder(),
                          prefixIcon: Icon(Icons.height),
                        ),
                        keyboardType: TextInputType.number,
                        validator: (value) {
                          if (value != null && value.trim().isNotEmpty) {
                            final height = double.tryParse(value.trim());
                            if (height == null || height < 50 || height > 300) {
                              return 'Height must be between 50-300 cm';
                            }
                          }
                          return null;
                        },
                      ),
                    ),
                    const SizedBox(width: 16),
                    Expanded(
                      child: TextFormField(
                        controller: _weightController,
                        decoration: const InputDecoration(
                          labelText: 'Weight (kg)',
                          border: OutlineInputBorder(),
                          prefixIcon: Icon(Icons.monitor_weight),
                        ),
                        keyboardType: TextInputType.number,
                        validator: (value) {
                          if (value != null && value.trim().isNotEmpty) {
                            final weight = double.tryParse(value.trim());
                            if (weight == null || weight < 20 || weight > 500) {
                              return 'Weight must be between 20-500 kg';
                            }
                          }
                          return null;
                        },
                      ),
                    ),
                  ],
                ),
                const SizedBox(height: 32),

                // Save Button
                SizedBox(
                  width: double.infinity,
                  child: BlocBuilder<UserCubit, UserState>(
                    builder: (context, state) {
                      return ElevatedButton(
                        onPressed: state.isUpdating ? null : _saveProfile,
                        style: ElevatedButton.styleFrom(
                          padding: const EdgeInsets.symmetric(vertical: 16),
                        ),
                        child: state.isUpdating
                            ? const SizedBox(
                                width: 20,
                                height: 20,
                                child: CircularProgressIndicator(
                                  strokeWidth: 2,
                                ),
                              )
                            : const Text('Save Changes'),
                      );
                    },
                  ),
                ),
                const SizedBox(height: 16),

                // Cancel Button
                SizedBox(
                  width: double.infinity,
                  child: OutlinedButton(
                    onPressed: () => context.pop(),
                    style: OutlinedButton.styleFrom(
                      padding: const EdgeInsets.symmetric(vertical: 16),
                    ),
                    child: const Text('Cancel'),
                  ),
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }

  Future<void> _selectDate() async {
    final DateTime? picked = await showDatePicker(
      context: context,
      initialDate:
          _selectedDate ??
          DateTime.now().subtract(const Duration(days: 365 * 20)),
      firstDate: DateTime(1900),
      lastDate: DateTime.now(),
    );
    if (picked != null && picked != _selectedDate) {
      setState(() {
        _selectedDate = picked;
      });
    }
  }

  void _saveProfile() {
    if (!_formKey.currentState!.validate()) {
      return;
    }

    final userCubit = context.read<UserCubit>();
    userCubit.updateProfile({
      'name': _nameController.text.trim(),
      'phone': _phoneController.text.trim().isEmpty
          ? null
          : _phoneController.text.trim(),
      'bio': _bioController.text.trim().isEmpty
          ? null
          : _bioController.text.trim(),
      'avatar': _avatarController.text.trim().isEmpty
          ? null
          : _avatarController.text.trim(),
      'date_of_birth': _selectedDate?.toIso8601String(),
      'sex': _selectedSex,
      'goal': _selectedGoal,
      'height_cm': _heightController.text.trim().isEmpty
          ? null
          : double.tryParse(_heightController.text.trim()),
      'weight_kg': _weightController.text.trim().isEmpty
          ? null
          : double.tryParse(_weightController.text.trim()),
    });
  }
}
