import 'package:dio/dio.dart';
import 'package:flutter/foundation.dart';
import '../network/api_helper.dart';
import '../network/dio_client.dart';

// Conditional import
import 'uploads_api_stub.dart'
    if (dart.library.html) 'uploads_api_web.dart'
    if (dart.library.io) 'uploads_api_io.dart';

class UploadsApi {
  final Dio _dio;
  UploadsApi(DioClient client) : _dio = client.dio;

  /// Upload an image file
  Future<Map<String, dynamic>> uploadImage(
    dynamic imageFile,
    String folder, {
    String? filename,
  }) async {
    try {
      MultipartFile multipartFile;

      if (kIsWeb) {
        final bytes = await readFileBytes(imageFile);
        final fileFilename = filename ?? _getFileName(imageFile);
        multipartFile = MultipartFile.fromBytes(bytes, filename: fileFilename);
      } else {
        final file = imageFile;
        multipartFile = await MultipartFile.fromFile(
          file.path,
          filename: filename ?? file.path.split('/').last,
        );
      }

      final formData = FormData.fromMap({
        'image': multipartFile,
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

  /// Helper to get file name (web-compatible)
  String _getFileName(dynamic file) {
    if (kIsWeb) {
      final timestamp = DateTime.now().millisecondsSinceEpoch;
      return 'upload_$timestamp.jpg';
    } else {
      final path = file.path;
      return path.split('/').last;
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
