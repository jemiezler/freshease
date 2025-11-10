// core/network/api_helper.dart
import 'package:dio/dio.dart';

/// Base helper class for consistent API error handling
class ApiHelper {
  /// Handle DioException and convert to user-friendly error messages
  static Exception handleDioError(DioException e) {
    if (e.response != null) {
      final statusCode = e.response?.statusCode;
      final message = e.response?.data is Map<String, dynamic>
          ? e.response?.data['message'] as String?
          : e.message;

      if (statusCode == 401) {
        return Exception('Authentication required. Please login again.');
      } else if (statusCode == 403) {
        return Exception('Access forbidden: ${message ?? 'Insufficient permissions'}');
      } else if (statusCode == 404) {
        return Exception('Resource not found');
      } else if (statusCode == 400) {
        return Exception('Invalid request: ${message ?? 'Please check your input'}');
      } else if (statusCode == 500) {
        return Exception('Server error: ${message ?? 'Please try again later'}');
      } else {
        return Exception('API Error: ${message ?? e.message}');
      }
    } else {
      return Exception('Network Error: ${e.message ?? 'Unable to connect to server'}');
    }
  }

  /// Extract data from response (handles both wrapped and direct responses)
  static T extractData<T>(dynamic responseData, T Function(Map<String, dynamic>) fromJson) {
    if (responseData is Map<String, dynamic>) {
      // Check if data is wrapped in 'data' field
      if (responseData.containsKey('data')) {
        final data = responseData['data'];
        if (data is Map<String, dynamic>) {
          return fromJson(data);
        } else if (data is List) {
          // For list responses, return the list directly
          return data as T;
        }
      }
      // Direct response
      return fromJson(responseData);
    }
    throw Exception('Invalid response format from server');
  }

  /// Extract list from response
  static List<T> extractList<T>(dynamic responseData, T Function(Map<String, dynamic>) fromJson) {
    if (responseData is Map<String, dynamic> && responseData.containsKey('data')) {
      final data = responseData['data'];
      if (data is List) {
        return data
            .map((item) => fromJson(item as Map<String, dynamic>))
            .toList();
      }
    } else if (responseData is List) {
      return responseData
          .map((item) => fromJson(item as Map<String, dynamic>))
          .toList();
    }
    throw Exception('Invalid response format from server');
  }
}

