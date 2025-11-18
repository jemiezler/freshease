// mobile implementation
import 'dart:io';
import 'dart:typed_data';

Future<List<int>> readFileBytes(File file) async {
  return await file.readAsBytes();
}
