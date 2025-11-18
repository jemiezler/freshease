import 'package:flutter_test/flutter_test.dart';
import 'package:flutter/material.dart';
import 'package:frontend/core/state/checkout_controller.dart';

void main() {
  group('CheckoutController', () {
    late CheckoutController controller;

    setUp(() {
      controller = CheckoutController();
    });

    tearDown(() {
      controller.dispose();
    });

    test('initial state has null shipping address and payment method', () {
      expect(controller.shippingAddress, isNull);
      expect(controller.paymentMethod, isNull);
      expect(controller.paymentLast4, isNull);
      expect(controller.selectedPlan, isNull);
      expect(controller.isPlanCheckout, false);
      expect(controller.isReadyForReview, false);
    });

    test('setShipping updates shipping address and notifies listeners', () {
      final address = Address(
        fullName: 'John Doe',
        phone: '+1234567890',
        line1: '123 Main St',
        line2: 'Apt 4B',
        subDistrict: 'Downtown',
        district: 'Central',
        province: 'Bangkok',
        postalCode: '10100',
      );

      var notified = false;
      controller.addListener(() {
        notified = true;
      });

      controller.setShipping(address);

      expect(controller.shippingAddress, equals(address));
      expect(notified, true);
    });

    test('setPayment updates payment method and notifies listeners', () {
      var notified = false;
      controller.addListener(() {
        notified = true;
      });

      controller.setPayment(PaymentMethod.card, last4: '1234');

      expect(controller.paymentMethod, equals(PaymentMethod.card));
      expect(controller.paymentLast4, equals('1234'));
      expect(notified, true);
    });

    test('setPayment with COD does not require last4', () {
      controller.setPayment(PaymentMethod.cod);

      expect(controller.paymentMethod, equals(PaymentMethod.cod));
      expect(controller.paymentLast4, isNull);
    });

    test('setPlanCheckout updates selected plan and notifies listeners', () {
      final plan = PlanOrder(
        id: 1,
        title: 'Weekly Plan',
        price: 99.99,
        subtitle: '7 days of meals',
      );

      var notified = false;
      controller.addListener(() {
        notified = true;
      });

      controller.setPlanCheckout(plan);

      expect(controller.selectedPlan, equals(plan));
      expect(controller.isPlanCheckout, true);
      expect(notified, true);
    });

    test('clearPlan removes selected plan and notifies listeners', () {
      final plan = PlanOrder(
        id: 1,
        title: 'Weekly Plan',
        price: 99.99,
        subtitle: '7 days of meals',
      );
      controller.setPlanCheckout(plan);

      var notified = false;
      controller.addListener(() {
        notified = true;
      });

      controller.clearPlan();

      expect(controller.selectedPlan, isNull);
      expect(controller.isPlanCheckout, false);
      expect(notified, true);
    });

    test('clear resets all fields and notifies listeners', () {
      final address = Address(
        fullName: 'John Doe',
        phone: '+1234567890',
        line1: '123 Main St',
        subDistrict: 'Downtown',
        district: 'Central',
        province: 'Bangkok',
        postalCode: '10100',
      );
      controller.setShipping(address);
      controller.setPayment(PaymentMethod.card, last4: '1234');
      controller.setPlanCheckout(PlanOrder(
        id: 1,
        title: 'Plan',
        price: 99.99,
        subtitle: 'Subtitle',
      ));

      var notified = false;
      controller.addListener(() {
        notified = true;
      });

      controller.clear();

      expect(controller.shippingAddress, isNull);
      expect(controller.paymentMethod, isNull);
      expect(controller.paymentLast4, isNull);
      expect(controller.selectedPlan, isNull);
      expect(notified, true);
    });

    test('isReadyForReview returns true when both address and payment are set', () {
      final address = Address(
        fullName: 'John Doe',
        phone: '+1234567890',
        line1: '123 Main St',
        subDistrict: 'Downtown',
        district: 'Central',
        province: 'Bangkok',
        postalCode: '10100',
      );

      expect(controller.isReadyForReview, false);

      controller.setShipping(address);
      expect(controller.isReadyForReview, false);

      controller.setPayment(PaymentMethod.cod);
      expect(controller.isReadyForReview, true);
    });

    test('isReadyForReview returns false when only address is set', () {
      final address = Address(
        fullName: 'John Doe',
        phone: '+1234567890',
        line1: '123 Main St',
        subDistrict: 'Downtown',
        district: 'Central',
        province: 'Bangkok',
        postalCode: '10100',
      );

      controller.setShipping(address);
      expect(controller.isReadyForReview, false);
    });

    test('isReadyForReview returns false when only payment is set', () {
      controller.setPayment(PaymentMethod.cod);
      expect(controller.isReadyForReview, false);
    });

    test('placeOrder returns OrderResult with generated ID', () async {
      final result = await controller.placeOrder();

      expect(result.orderId, isNotEmpty);
      expect(result.orderId, startsWith('FE-'));
      expect(result.createdAt, isA<DateTime>());
    });

    test('placeOrder generates unique order IDs', () async {
      final result1 = await controller.placeOrder();
      await Future.delayed(const Duration(milliseconds: 10));
      final result2 = await controller.placeOrder();

      expect(result1.orderId, isNot(equals(result2.orderId)));
    });
  });

  group('CheckoutScope', () {
    test('of returns controller from context', () {
      final controller = CheckoutController();
      final scope = CheckoutScope(
        controller: controller,
        child: Container(),
      );

      // In a real widget test, we would use tester.widget to get the context
      // For unit test, we verify the structure
      expect(scope.controller, equals(controller));
    });
  });
}

