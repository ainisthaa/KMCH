import React from "react";
import MultiSelectDropdown from "./MultiSelectDropdown";


function FoodModal({
  show,
  isEdit,
  isDelete,
  onClose,
  onSubmit,
  onDelete,
  name,
  setName,
  price,
  setPrice,
  categoryId,
  setCategoryId,
  categoryList,
  ingredientList,
  selectedIngredientIds,
  setSelectedIngredientIds,
}) {
  if (!show) return null;

  return (
    
    <div className="fixed inset-0 flex items-center justify-center z-50 pointer-events-none">
      <div className="bg-white p-6 rounded-xl w-[300px] text-center shadow-lg pointer-events-auto">
        
        
        {isDelete ? (
          <>
            <p className="mb-4 font-bold text-rose-700">ต้องการลบรายการนี้หรือไม่?</p>
            <div className="flex justify-between">
              <button
                className="bg-blue-500 text-white px-3 py-1 rounded"
                onClick={onDelete}
              >
                ยืนยัน
              </button>
              <button
                className="bg-gray-300 text-black px-3 py-1 rounded"
                onClick={onClose}
              >
                ยกเลิก
              </button>
            </div>
          </>
        ) : (
          <>
            
            <h3 className="text-lg font-bold text-rose-600 mb-3">
              {isEdit ? "แก้ไขอาหาร" : "เพิ่มอาหาร"}
            </h3>
            <p className="text-left font-semibold">อาหาร:</p>
            <input
              type="text"
              className="border w-full p-2 mb-3"
              placeholder="ชื่ออาหาร"
              value={name}
              onChange={(e) => setName(e.target.value)}
            />
            <p className="text-left font-semibold">ราคา:</p>
            <input
              type="number"
              className="border w-full p-2 mb-3"
              placeholder="ราคา"
              value={price}
              onChange={(e) => setPrice(e.target.value)}
            />
            <p className="text-left font-semibold">categoryid:</p>
            <select
              value={categoryId}
              onChange={(e) => setCategoryId(Number(e.target.value))}
              className="border w-full p-2 mb-3"
            >
              <option value="">-- เลือกหมวดหมู่ --</option>
              {categoryList.map((cat) => (
                <option key={cat.id} value={cat.id}>
                  {cat.name}
                </option>
              ))}
            </select>
            <div className="text-left mb-3 w-full">
            <p className="font-semibold mb-1">วัตถุดิบ:</p>
            <MultiSelectDropdown
                options={ingredientList}
                selected={selectedIngredientIds}
                setSelected={setSelectedIngredientIds}
            />
            </div>


            <div className="flex justify-between">
              <button
                className="bg-blue-500 text-white px-3 py-1 rounded"
                onClick={onSubmit}
              >
                Submit
              </button>
              <button
                className="bg-gray-300 text-black px-3 py-1 rounded"
                onClick={onClose}
              >
                Cancel
              </button>
            </div>
          </>
        )}
      </div>
    </div>
  );
}

export default FoodModal;
