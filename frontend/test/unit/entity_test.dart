import 'package:flutter_test/flutter_test.dart';
import 'package:frontend/features/account/domain/entities/user_profile.dart';
import 'package:frontend/features/auth/domain/entities/user.dart';

void main() {
  group('UserProfile', () {
    test('initials returns first letter for single name', () {
      final profile = UserProfile(
        id: '1',
        email: 'test@example.com',
        name: 'John',
        status: 'active',
        createdAt: DateTime.now(),
        updatedAt: DateTime.now(),
      );

      expect(profile.initials, 'J');
    });

    test('initials returns first letters for full name', () {
      final profile = UserProfile(
        id: '1',
        email: 'test@example.com',
        name: 'John Doe',
        status: 'active',
        createdAt: DateTime.now(),
        updatedAt: DateTime.now(),
      );

      expect(profile.initials, 'JD');
    });

    test('initials returns first letter for single character name', () {
      final profile = UserProfile(
        id: '1',
        email: 'test@example.com',
        name: 'J',
        status: 'active',
        createdAt: DateTime.now(),
        updatedAt: DateTime.now(),
      );

      expect(profile.initials, 'J');
    });

    test('initials handles normal name with space', () {
      final profile = UserProfile(
        id: '1',
        email: 'test@example.com',
        name: 'John Doe',
        status: 'active',
        createdAt: DateTime.now(),
        updatedAt: DateTime.now(),
      );

      expect(profile.initials, 'JD');
    });

    test('displayName returns name when not empty', () {
      final profile = UserProfile(
        id: '1',
        email: 'test@example.com',
        name: 'John Doe',
        status: 'active',
        createdAt: DateTime.now(),
        updatedAt: DateTime.now(),
      );

      expect(profile.displayName, 'John Doe');
    });

    test('displayName returns email when name is empty', () {
      final profile = UserProfile(
        id: '1',
        email: 'test@example.com',
        name: '',
        status: 'active',
        createdAt: DateTime.now(),
        updatedAt: DateTime.now(),
      );

      expect(profile.displayName, 'test@example.com');
    });

    test('isCompleteProfile returns true when all required fields present', () {
      final profile = UserProfile(
        id: '1',
        email: 'test@example.com',
        name: 'John Doe',
        phone: '+1234567890',
        status: 'active',
        createdAt: DateTime.now(),
        updatedAt: DateTime.now(),
      );

      expect(profile.isCompleteProfile, true);
    });

    test('isCompleteProfile returns false when phone is missing', () {
      final profile = UserProfile(
        id: '1',
        email: 'test@example.com',
        name: 'John Doe',
        status: 'active',
        createdAt: DateTime.now(),
        updatedAt: DateTime.now(),
      );

      expect(profile.isCompleteProfile, false);
    });

    test('isCompleteProfile returns false when name is empty', () {
      final profile = UserProfile(
        id: '1',
        email: 'test@example.com',
        name: '',
        phone: '+1234567890',
        status: 'active',
        createdAt: DateTime.now(),
        updatedAt: DateTime.now(),
      );

      expect(profile.isCompleteProfile, false);
    });

    test('isCompleteProfile returns false when email is empty', () {
      final profile = UserProfile(
        id: '1',
        email: '',
        name: 'John Doe',
        phone: '+1234567890',
        status: 'active',
        createdAt: DateTime.now(),
        updatedAt: DateTime.now(),
      );

      expect(profile.isCompleteProfile, false);
    });

    test('copyWith creates new instance with updated fields', () {
      final original = UserProfile(
        id: '1',
        email: 'test@example.com',
        name: 'John Doe',
        status: 'active',
        createdAt: DateTime.parse('2024-01-01T00:00:00Z'),
        updatedAt: DateTime.parse('2024-01-01T00:00:00Z'),
      );

      final updated = original.copyWith(
        name: 'Jane Doe',
        phone: '+9876543210',
      );

      expect(updated.id, '1');
      expect(updated.email, 'test@example.com');
      expect(updated.name, 'Jane Doe');
      expect(updated.phone, '+9876543210');
      expect(updated.status, 'active');
    });

    test('copyWith preserves original values when fields not provided', () {
      final original = UserProfile(
        id: '1',
        email: 'test@example.com',
        name: 'John Doe',
        status: 'active',
        createdAt: DateTime.parse('2024-01-01T00:00:00Z'),
        updatedAt: DateTime.parse('2024-01-01T00:00:00Z'),
      );

      final updated = original.copyWith();

      expect(updated.id, original.id);
      expect(updated.email, original.email);
      expect(updated.name, original.name);
      expect(updated.status, original.status);
    });

    test('equality works correctly', () {
      final profile1 = UserProfile(
        id: '1',
        email: 'test@example.com',
        name: 'John Doe',
        status: 'active',
        createdAt: DateTime.parse('2024-01-01T00:00:00Z'),
        updatedAt: DateTime.parse('2024-01-01T00:00:00Z'),
      );

      final profile2 = UserProfile(
        id: '1',
        email: 'test@example.com',
        name: 'John Doe',
        status: 'active',
        createdAt: DateTime.parse('2024-01-01T00:00:00Z'),
        updatedAt: DateTime.parse('2024-01-01T00:00:00Z'),
      );

      expect(profile1, equals(profile2));
    });
  });

  group('User', () {
    test('equality works correctly', () {
      final user1 = User(
        id: '1',
        email: 'test@example.com',
        name: 'John Doe',
        avatar: 'avatar.jpg',
      );

      final user2 = User(
        id: '1',
        email: 'test@example.com',
        name: 'John Doe',
        avatar: 'avatar.jpg',
      );

      expect(user1, equals(user2));
    });

    test('equality returns false for different users', () {
      final user1 = User(
        id: '1',
        email: 'test@example.com',
        name: 'John Doe',
      );

      final user2 = User(
        id: '2',
        email: 'test2@example.com',
        name: 'Jane Doe',
      );

      expect(user1, isNot(equals(user2)));
    });

    test('handles null optional fields', () {
      final user = User(
        id: '1',
        email: 'test@example.com',
      );

      expect(user.name, isNull);
      expect(user.avatar, isNull);
    });
  });
}

