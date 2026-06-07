// ปุ่มเพิ่ม
export function AddButton({ onClick }) {
  return (
    <button
      className="bg-rose-400 text-white px-4 py-2 rounded font-bold"
      onClick={onClick}
    >
      add
    </button>
  );
}

// ปุ่มแก้ไข
export function EditButton({ onClick }) {
  return (
    <button
      className="bg-green-500 text-white px-2 py-0.5 rounded"
      onClick={onClick}
    >
      edit
    </button>
  );
}

// ปุ่มลบ
export function DeleteButton({ onClick }) {
  return (
    <button
      className="bg-red-400 text-white px-2 py-0.5 rounded ml-2"
      onClick={onClick}
    >
      ❌
    </button>
  );
}

export function ViewButton({ onClick }) {
  return (
    <button
      className="bg-purple-400 text-white px-3 py-1 rounded ml-2"
      onClick={onClick}
    >
      view
    </button>
  );
}
