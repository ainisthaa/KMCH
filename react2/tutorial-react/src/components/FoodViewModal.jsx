import React from "react";

function FoodViewModal({ foodDetail, onClose, categoryList }) {
  if (!foodDetail) return null;
  console.log(categoryList)
  const categoryFound = categoryList.find((cat) => {
    if (cat.id === foodDetail.category_id) {
      return cat
    }
  })
  // console.log("Test", categoryFound.name)
  return (
    <div className="fixed inset-0 flex items-center justify-center z-50 bg-opacity-30">
      <div className="bg-white p-5 rounded-xl w-[300px] text-left shadow-md">
        <h3 className="font-bold text-lg mb-3 text-rose-700">รายละเอียดอาหาร</h3>

        <p><span className="font-semibold">ชื่อ:</span> {foodDetail.name}</p>
        <p><span className="font-semibold">ราคา:</span> {foodDetail.price}</p>
        <p><span className="font-semibold">หมวดหมู่:</span> {categoryFound.name}</p>

        <div className="mt-3">
          <p className="font-semibold">วัตถุดิบ:</p>
          {foodDetail.ingredients.length === 0 ? (
            <p className="text-gray-500">ไม่มีวัตถุดิบ</p>
          ) : (
            <ul className="text-sm list-disc list-inside">
              {foodDetail.ingredients.map((i) => (
                <li key={i.ingredient_id}>{i.ingredient_name}</li>
              ))}
            </ul>
          )}
        </div>

        <button
          className="mt-4 bg-gray-300 px-4 py-1 rounded"
          onClick={onClose}
        >
          ปิด
        </button>
      </div>
    </div>
  );
}

export default FoodViewModal;
