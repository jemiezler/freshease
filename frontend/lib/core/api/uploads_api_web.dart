// web implementation
import 'dart:html' as html;
import 'dart:typed_data';
import 'dart:async';

Future<List<int>> readFileBytes(html.File file) {
  final completer = Completer<List<int>>();
  final reader = html.FileReader();

  reader.onLoadEnd.listen((event) {
    final result = reader.result;
    if (result != null) {
      final bytes = Uint8List.view(result as ByteBuffer);
      completer.complete(bytes);
    } else {
      completer.completeError('Failed to read file');
    }
  });

  reader.onError.listen((event) {
    completer.completeError('File read error');
  });

  reader.readAsArrayBuffer(file);

  return completer.future;
}
