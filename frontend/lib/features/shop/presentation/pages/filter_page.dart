// ไฟล์: filter_page.dart
import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

// ธีมสี
const Color _primaryColor = Color(0xFF90B56D);

class FilterPage extends StatefulWidget {
  // รับค่า filters ปัจจุบันเข้ามา
  final Map<String, List<String>> initialFilters;

  const FilterPage({super.key, required this.initialFilters});

  @override
  State<FilterPage> createState() => _FilterPageState();
}

class _FilterPageState extends State<FilterPage> {
  // สถานะปัจจุบันของฟิลเตอร์
  late List<String> _selectedCategories;
  late List<String> _selectedBrands;

  // ข้อมูลจำลองสำหรับฟิลเตอร์
  final List<String> _allCategories = const [
    'Eggs',
    'Noodles & Pasta',
    'Chips & Crisps',
    'Fast Food',
  ];
  final List<String> _allBrands = const [
    'Individual Collection',
    'Malee',
    'Ifod',
    'Kazi Farmas',
  ];

  @override
  void initState() {
    super.initState();
    // คัดลอกค่าเริ่มต้นมาใส่ใน state
    _selectedCategories = List.from(widget.initialFilters['categories'] ?? []);
    _selectedBrands = List.from(widget.initialFilters['brands'] ?? []);
  }

  // ฟังก์ชันสำหรับ Apply และส่งค่ากลับ
  void _applyFilters() {
    final result = {
      'categories': _selectedCategories,
      'brands': _selectedBrands,
    };
    // ส่งค่ากลับไปหน้า ShopPage
    context.pop(result);
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: Colors.white,
      appBar: AppBar(
        backgroundColor: Colors.white,
        elevation: 0,
        centerTitle: true,
        title: const Text(
          'Filters',
          style: TextStyle(color: Colors.black, fontWeight: FontWeight.bold),
        ),
        leading: IconButton(
          icon: const Icon(Icons.close, color: Colors.black),
          onPressed: () => context.pop(), // กดปิดโดยไม่ Apply
        ),
      ),
      body: ListView(
        padding: const EdgeInsets.all(24),
        children: [
          _buildSectionTitle('Categories'),
          ..._allCategories.map(
            (category) => _buildCheckbox(
              title: category,
              selectedList: _selectedCategories,
            ),
          ),

          const SizedBox(height: 24),

          _buildSectionTitle('Brand'),
          ..._allBrands.map(
            (brand) =>
                _buildCheckbox(title: brand, selectedList: _selectedBrands),
          ),
        ],
      ),
      bottomNavigationBar: Padding(
        padding: EdgeInsets.fromLTRB(
          24,
          16,
          24,
          MediaQuery.of(context).padding.bottom + 16,
        ),
        child: ElevatedButton(
          onPressed: _applyFilters,
          style: ElevatedButton.styleFrom(
            backgroundColor: _primaryColor,
            minimumSize: const Size(double.infinity, 56),
            shape: RoundedRectangleBorder(
              borderRadius: BorderRadius.circular(16),
            ),
            elevation: 0,
          ),
          child: const Text(
            'Apply Filter',
            style: TextStyle(
              fontSize: 18,
              color: Colors.white,
              fontWeight: FontWeight.bold,
            ),
          ),
        ),
      ),
    );
  }

  // --- Widget สำหรับสร้าง CheckboxListTile ---
  Widget _buildCheckbox({
    required String title,
    required List<String> selectedList,
  }) {
    final isSelected = selectedList.contains(title);
    return CheckboxListTile(
      value: isSelected,
      title: Text(title),
      activeColor: _primaryColor,
      controlAffinity: ListTileControlAffinity.leading, // Checkbox อยู่ด้านซ้าย
      onChanged: (bool? value) {
        setState(() {
          if (value == true) {
            selectedList.add(title);
          } else {
            selectedList.remove(title);
          }
        });
      },
    );
  }

  // --- Widget สำหรับสร้าง Title ของ Section ---
  Widget _buildSectionTitle(String title) {
    return Padding(
      padding: const EdgeInsets.only(bottom: 8.0),
      child: Text(
        title,
        style: const TextStyle(fontSize: 20, fontWeight: FontWeight.bold),
      ),
    );
  }
}
