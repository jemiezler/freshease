// ignore_for_file: constant_identifier_names
// ignore: deprecated_member_use
import 'package:flutter/material.dart';

// ⚠️ ต้องแน่ใจว่า AppColors_primary ถูก import หรือนิยามไว้
const Color AppColors_primary = Color(0xFF90B56D);

class FilterBottomSheet extends StatefulWidget {
  final List<String> initialCategories;
  final String? initialBrand;
  final Function(List<String> categories, String? brand) onApplyFilter;

  const FilterBottomSheet({
    super.key,
    required this.initialCategories,
    required this.initialBrand,
    required this.onApplyFilter,
  });

  @override
  State<FilterBottomSheet> createState() => _FilterBottomSheetState();
}

class _FilterBottomSheetState extends State<FilterBottomSheet> {
  // 🎯 ‼️ แก้ไขตรงนี้: เปลี่ยน Categories ให้ตรงกับ explore_page.dart
  final List<String> _allCategories = [
    'Fruits',
    'Oil',
    'Meat',
    'Bakery',
    'Dairy',
    'Beverages',
  ];

  // 🎯 ข้อมูล Brand (ยังคงเดิม)
  final List<String> _allBrands = [
    'Individual Collection',
    'Malee',
    'Ifod',
    'Kazi Farmas',
  ];

  late List<String> _currentCategories;
  String? _currentBrand;

  @override
  void initState() {
    super.initState();
    _currentCategories = List.from(widget.initialCategories);
    _currentBrand = widget.initialBrand;
  }

  Widget _buildHeader(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.only(left: 16, right: 16, top: 10),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [
          IconButton(
            icon: const Icon(Icons.close, size: 28),
            onPressed: () => Navigator.pop(context),
          ),
          Text(
            'Filters',
            style: Theme.of(context).textTheme.titleLarge?.copyWith(
              fontWeight: FontWeight.w700,
              fontSize: 22,
            ),
          ),
          const SizedBox(width: 48),
        ],
      ),
    );
  }

  Widget _buildFilterSection({
    required String title,
    required List<Widget> children,
  }) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Padding(
          padding: const EdgeInsets.only(left: 20, top: 20, bottom: 8),
          child: Text(
            title,
            style: Theme.of(
              context,
            ).textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold),
          ),
        ),
        ...children,
        const SizedBox(height: 10),
        const Divider(
          height: 1,
          thickness: 1,
          indent: 16,
          endIndent: 16,
          color: Colors.black12,
        ),
      ],
    );
  }

  @override
  Widget build(BuildContext context) {
    return DraggableScrollableSheet(
      initialChildSize: 0.95,
      minChildSize: 0.5,
      maxChildSize: 0.95,
      expand: false,
      builder: (context, scrollController) {
        return Container(
          decoration: const BoxDecoration(
            color: Colors.white,
            borderRadius: BorderRadius.vertical(top: Radius.circular(20)),
          ),
          child: Column(
            children: [
              _buildHeader(context),
              Expanded(
                child: ListView(
                  controller: scrollController,
                  padding: EdgeInsets.zero,
                  children: [
                    // --- Categories Filter (Multiple Select) ---
                    _buildFilterSection(
                      title: 'Categories',
                      children: _allCategories.map((category) {
                        final isSelected = _currentCategories.contains(
                          category,
                        );
                        return Padding(
                          padding: const EdgeInsets.symmetric(horizontal: 4.0),
                          child: CheckboxListTile(
                            dense: true,
                            title: Text(category), // 👈 จะแสดง Fruits, Oil...
                            value: isSelected,
                            onChanged: (bool? newValue) {
                              setState(() {
                                if (newValue == true) {
                                  _currentCategories.add(category);
                                } else {
                                  _currentCategories.remove(category);
                                }
                              });
                            },
                            controlAffinity: ListTileControlAffinity.leading,
                            checkColor: Colors.white,
                            activeColor: AppColors_primary,
                          ),
                        );
                      }).toList(),
                    ),
                    // --- Brand Filter (Single Select) ---
                    _buildFilterSection(
                      title: 'Brand',
                      children: _allBrands.map((brand) {
                        return Padding(
                          padding: const EdgeInsets.symmetric(horizontal: 4.0),
                          child: RadioListTile<String>(
                            dense: true,
                            title: Text(brand),
                            value: brand,
                            // ignore: deprecated_member_use
                            groupValue: _currentBrand,
                            // ignore: deprecated_member_use
                            onChanged: (String? newValue) {
                              setState(() {
                                _currentBrand = newValue;
                              });
                            },
                            activeColor: AppColors_primary,
                            controlAffinity: ListTileControlAffinity.leading,
                          ),
                        );
                      }).toList(),
                    ),
                    const SizedBox(height: 100),
                  ],
                ),
              ),
              // --- Apply Filter Button (Bottom) ---
              Padding(
                padding: const EdgeInsets.only(
                  left: 20,
                  right: 20,
                  bottom: 30,
                  top: 10,
                ),
                child: SizedBox(
                  width: double.infinity,
                  height: 60,
                  child: ElevatedButton(
                    onPressed: () {
                      // 🎯 ส่งค่า ['Fruits', 'Oil'] กลับไป
                      widget.onApplyFilter(_currentCategories, _currentBrand);
                    },
                    style: ElevatedButton.styleFrom(
                      backgroundColor: AppColors_primary,
                      shape: RoundedRectangleBorder(
                        borderRadius: BorderRadius.circular(20),
                      ),
                      elevation: 0,
                    ),
                    child: Text(
                      'Apply Filter',
                      style: Theme.of(context).textTheme.titleMedium?.copyWith(
                        color: Colors.white,
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                  ),
                ),
              ),
            ],
          ),
        );
      },
    );
  }
}
