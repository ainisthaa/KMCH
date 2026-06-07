function SortDropdown({ sortKey, setSortKey, setPage }) {
    return (
      <div>
        <p>เรียงตาม</p>
        <select
          className="bg-amber-200"
          value={sortKey}
          onChange={(e) => {
            setSortKey(e.target.value);
            setPage(1); // reset หน้า
          }}
        >
          <option value="">-- Default (ID) --</option>
          <option value="name">ชื่อ A-Z</option>
          <option value="-name">ชื่อ Z-A</option>
          <option value="price">ราคาน้อย → มาก</option>
          <option value="-price">ราคามาก → น้อย</option>
          <option value="category">หมวดหมู่ (น้อย → มาก)</option>
          <option value="-category">หมวดหมู่ (มาก → น้อย)</option>
          <option value="updated">อัปเดตเก่า → ใหม่</option>
          <option value="-updated">อัปเดตใหม่ → เก่า</option>
        </select>
      </div>
    );
  }
  
  export default SortDropdown;
  