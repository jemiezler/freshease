// ไฟล์: search_pill_widget.dart

import 'package:flutter/material.dart';

class SearchPill extends StatelessWidget {
  final TextEditingController? controller;
  final VoidCallback? onFilterTap;
  final VoidCallback? onTap;
  final ValueChanged<String>? onChanged;
  final bool showFilter;
  final bool readOnly;

  const SearchPill({
    super.key,
    this.controller,
    this.onFilterTap,
    this.onTap,
    this.onChanged,
    this.showFilter = false,
    this.readOnly = false,
  });

  @override
  Widget build(BuildContext context) {
    return Container(
      height: 48,
      decoration: BoxDecoration(
        // ‼️ แก้ไข: เปลี่ยนสีพื้นหลังเป็นเทาอ่อน
        color: const Color(0xFFF2F3F2),
        borderRadius: BorderRadius.circular(24),
        // ‼️ ลบ border และ boxShadow ออก
      ),
      padding: const EdgeInsets.symmetric(horizontal: 14),
      child: Row(
        children: [
          Icon(Icons.search, size: 22, color: Colors.grey.shade600),
          const SizedBox(width: 8),
          Expanded(
            child: TextField(
              controller: controller,
              readOnly: readOnly,
              onTap: onTap,
              onChanged: onChanged,
              decoration: InputDecoration(
                border: InputBorder.none,
                hintText: 'Search Store',
                // ‼️ แก้ไข: สี Hint ให้เข้มขึ้นเล็กน้อย
                hintStyle: TextStyle(color: Colors.grey.shade600),
              ),
              textInputAction: TextInputAction.search,
            ),
          ),
          if (showFilter)
            IconButton(
              icon: Icon(Icons.tune_rounded, color: Colors.grey.shade600),
              onPressed: onFilterTap,
            ),
        ],
      ),
    );
  }
}
