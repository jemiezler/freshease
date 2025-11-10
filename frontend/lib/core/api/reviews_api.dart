// core/api/reviews_api.dart
import 'package:dio/dio.dart';
import '../network/api_helper.dart';
import '../network/dio_client.dart';

class ReviewsApi {
  final Dio _dio;
  ReviewsApi(DioClient client) : _dio = client.dio;

  /// Get all reviews
  Future<List<Map<String, dynamic>>> listReviews() async {
    try {
      final response = await _dio.get('/api/reviews');
      if (response.statusCode == 200) {
        return ApiHelper.extractList<Map<String, dynamic>>(
          response.data,
          (json) => json,
        );
      }
      throw Exception('Failed to fetch reviews: ${response.statusCode}');
    } on DioException catch (e) {
      throw ApiHelper.handleDioError(e);
    }
  }

  /// Get review by ID
  Future<Map<String, dynamic>> getReview(String reviewId) async {
    try {
      final response = await _dio.get('/api/reviews/$reviewId');
      if (response.statusCode == 200) {
        return ApiHelper.extractData<Map<String, dynamic>>(
          response.data,
          (json) => json,
        );
      }
      throw Exception('Failed to fetch review: ${response.statusCode}');
    } on DioException catch (e) {
      throw ApiHelper.handleDioError(e);
    }
  }

  /// Create a new review
  Future<Map<String, dynamic>> createReview(Map<String, dynamic> reviewData) async {
    try {
      final response = await _dio.post('/api/reviews', data: reviewData);
      if (response.statusCode == 201 || response.statusCode == 200) {
        return ApiHelper.extractData<Map<String, dynamic>>(
          response.data,
          (json) => json,
        );
      }
      throw Exception('Failed to create review: ${response.statusCode}');
    } on DioException catch (e) {
      throw ApiHelper.handleDioError(e);
    }
  }

  /// Update a review
  Future<Map<String, dynamic>> updateReview(String reviewId, Map<String, dynamic> reviewData) async {
    try {
      final response = await _dio.patch('/api/reviews/$reviewId', data: reviewData);
      if (response.statusCode == 200 || response.statusCode == 201) {
        return ApiHelper.extractData<Map<String, dynamic>>(
          response.data,
          (json) => json,
        );
      }
      throw Exception('Failed to update review: ${response.statusCode}');
    } on DioException catch (e) {
      throw ApiHelper.handleDioError(e);
    }
  }

  /// Delete a review
  Future<void> deleteReview(String reviewId) async {
    try {
      final response = await _dio.delete('/api/reviews/$reviewId');
      if (response.statusCode != 202 && response.statusCode != 200) {
        throw Exception('Failed to delete review: ${response.statusCode}');
      }
    } on DioException catch (e) {
      throw ApiHelper.handleDioError(e);
    }
  }
}

