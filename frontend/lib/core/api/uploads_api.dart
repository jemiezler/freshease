// core/api/uploads_api.dart
import 'package:dio/dio.dart';
import 'dart:io';
import '../network/api_helper.dart';
import '../network/dio_client.dart';

class UploadsApi {
  final Dio _dio;
  UploadsApi(DioClient client) : _dio = client.dio;

  /// Upload an image file
  /// Returns the object name and URL
  Future<Map<String, dynamic>> uploadImage(File imageFile, String folder) async {
    try {
      final formData = FormData.fromMap({
        'image': await MultipartFile.fromFile(
          imageFile.path,
          filename: imageFile.path.split('/').last,
        ),
        'folder': folder,
      });

      final response = await _dio.post('/api/uploads/image', data: formData);
      if (response.statusCode == 200 || response.statusCode == 201) {
        return ApiHelper.extractData<Map<String, dynamic>>(
          response.data,
          (json) => json,
        );
      }
      throw Exception('Failed to upload image: ${response.statusCode}');
    } on DioException catch (e) {
      throw ApiHelper.handleDioError(e);
    }
  }

  /// Get image URL by object name
  Future<String> getImageUrl(String objectName) async {
    try {
      final response = await _dio.get('/api/uploads/image/$objectName');
      if (response.statusCode == 200) {
        final data = response.data;
        if (data is Map<String, dynamic> && data.containsKey('url')) {
          return data['url'] as String;
        } else if (data is String) {
          return data;
        }
        throw Exception('Invalid response format');
      }
      throw Exception('Failed to get image URL: ${response.statusCode}');
    } on DioException catch (e) {
      throw ApiHelper.handleDioError(e);
    }
  }
}

