// core/api/uploads_api.dart
import 'package:dio/dio.dart';
import 'package:flutter/foundation.dart';
import 'dart:async';
import '../network/api_helper.dart';
import '../network/dio_client.dart';

// Conditional imports for file handling
import 'dart:io' if (dart.library.html) 'dart:html' as io;
import 'dart:typed_data';

class UploadsApi {
  final Dio _dio;
  UploadsApi(DioClient client) : _dio = client.dio;

  /// Upload an image file
  /// Returns the object name and URL
  /// Accepts both File (mobile) and web File types
  /// For web, optionally provide a filename parameter
  Future<Map<String, dynamic>> uploadImage(
    dynamic imageFile,
    String folder, {
    String? filename,
  }) async {
    try {
      MultipartFile multipartFile;

      if (kIsWeb) {
        // Web: imageFile should be html.File or have bytes
        // For web, we expect the file to have a readAsBytes method
        final bytes = await _readFileBytes(imageFile);
        final fileFilename = filename ?? _getFileName(imageFile);
        multipartFile = MultipartFile.fromBytes(bytes, filename: fileFilename);
      } else {
        // Mobile: imageFile should be File
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

  /// Helper to read file bytes (web-compatible)
  Future<List<int>> _readFileBytes(dynamic file) async {
    if (kIsWeb) {
      // Web: html.File - use FileReader to read as ArrayBuffer
      // Note: For web, files are typically already available as bytes
      // This implementation handles html.File from file input elements
      return _readWebFile(file);
    } else {
      // Mobile: file is io.File
      return await file.readAsBytes();
    }
  }

  /// Read web file using FileReader API
  Future<List<int>> _readWebFile(dynamic file) async {
    final completer = Completer<List<int>>();

    // Use dynamic to work around conditional import type issues
    // FileReader is available in dart:html when compiling for web
    try {
      // Access FileReader through the io namespace (which is dart:html on web)
      // We need to use a workaround since FileReader constructor isn't directly accessible
      // with conditional imports in the analyzer
      final fileReader = _createFileReader();
      final htmlFile = file as io.File;

      fileReader.onLoadEnd.listen((e) {
        final result = fileReader.result;
        if (result != null) {
          try {
            final arrayBuffer = result as dynamic;
            final bytes = Uint8List.view(arrayBuffer);
            completer.complete(bytes);
          } catch (e) {
            completer.completeError('Failed to convert file: $e');
          }
        } else {
          completer.completeError('Failed to read file - no result');
        }
      });

      fileReader.onError.listen((e) {
        completer.completeError('File read error: $e');
      });

      fileReader.readAsArrayBuffer(htmlFile);
    } catch (e) {
      completer.completeError('Failed to create FileReader: $e');
    }

    return completer.future;
  }

  /// Create FileReader instance (web only)
  dynamic _createFileReader() {
    // Use noSuchMethod or reflection to create FileReader
    // Since we can't directly reference FileReader with conditional imports,
    // we'll use a workaround by accessing it through the html library
    if (kIsWeb) {
      // On web, FileReader is available in dart:html
      // We access it through the io namespace which is dart:html on web
      // ignore: undefined_function
      return io.FileReader();
    }
    throw UnsupportedError('FileReader only available on web');
  }

  /// Helper to get file name (web-compatible)
  String _getFileName(dynamic file) {
    if (kIsWeb) {
      // Web: html.File doesn't have a direct name property
      // The name typically comes from the FileList/InputElement
      // For now, use a default filename with timestamp
      final timestamp = DateTime.now().millisecondsSinceEpoch;
      return 'upload_$timestamp.jpg';
    } else {
      // Mobile: io.File has path
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
