import React from "react";

function ModalComponent({
  showModal,
  setShowModal,
  showAddModal,
  setShowAddModal,
  showDeleteModal,
  setShowDeleteModal,
  editValue,
  setEditValue,
  onEditSubmit,
  addValue,
  setAddValue,
  onAddSubmit,
  onDelete
}) {
  if (!showModal && !showAddModal && !showDeleteModal) return null;

  return (
    <div className="fixed inset-0 flex items-center justify-center z-50">
      <div className="bg-white p-6 rounded-xl w-[300px] text-center shadow-2xl">
        {/* Modal แก้ไข */}
        {showModal && (
          <>

            <h3 className="text-lg font-bold text-rose-600 mb-3">แก้ไขรายการ</h3>
            <p className="text-left font-semibold">อาหาร:</p>
            <input
              type="text"
              value={editValue}
              onChange={(e) => setEditValue(e.target.value)}
              className="border w-full p-2 mb-3"
              placeholder="ชื่อรายการใหม่"
            />
            <div className="flex justify-between">
              <button className="bg-blue-500 text-white px-3 py-1 rounded" onClick={onEditSubmit}>
                บันทึก
              </button>
              <button
                className="bg-gray-300 px-3 py-1 rounded"
                onClick={() => setShowModal(false)}
              >
                ยกเลิก
              </button>
            </div>
          </>
        )}

        {/* Modal เพิ่มรายการ */}
        {showAddModal && (
          <>
            <h3 className="text-lg font-bold text-rose-600 mb-3">เพิ่มรายการใหม่</h3>
            <p className="text-left font-semibold">อาหาร:</p>
            <input
              type="text"
              value={addValue}
              onChange={(e) => setAddValue(e.target.value)}
              className="border w-full p-2 mb-3"
              placeholder="ชื่อรายการ"
            />
            <div className="flex justify-between">
              <button className="bg-green-600 text-white px-3 py-1 rounded" onClick={onAddSubmit}>
                เพิ่ม
              </button>
              <button
                className="bg-gray-300 px-3 py-1 rounded"
                onClick={() => setShowAddModal(false)}
              >
                ยกเลิก
              </button>
            </div>
          </>
        )}

        {/* Modal ลบ */}
        {showDeleteModal && (
          <>
            <h3 className="text-lg font-bold text-red-600 mb-4">ต้องการลบรายการนี้หรือไม่?</h3>
            <div className="flex justify-between">
              <button className="bg-red-500 text-white px-3 py-1 rounded" onClick={onDelete}>
                ยืนยัน
              </button>
              <button
                className="bg-gray-300 px-3 py-1 rounded"
                onClick={() => setShowDeleteModal(false)}
              >
                ยกเลิก
              </button>
            </div>
          </>
        )}
      </div>
    </div>
  );
}

export default ModalComponent;
