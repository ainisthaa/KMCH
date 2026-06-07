import React, { useState } from "react";

function MultiSelectDropdown({ options, selected, setSelected }) {
  const [open, setOpen] = useState(false);

  const toggleSelect = (id) => {
    if (selected.includes(id)) {
      setSelected(selected.filter((s) => s !== id));
    } else {
      setSelected([...selected, id]);
    }
  };

  return (
    <div className="relative">
      <div
        className="border p-2 rounded bg-white cursor-pointer"
        onClick={() => setOpen(!open)}
      >
        {selected.length > 0
          ? options
              .filter((opt) => selected.includes(opt.id))
              .map((opt) => opt.name)
              .join(", ")
          : "เลือกวัตถุดิบ"}
      </div>

      {open && (
        <div className="absolute z-10 bg-white border w-full mt-1 max-h-40 overflow-y-auto rounded shadow">
          {options.map((option) => (
            <label
              key={option.id}
              className="block px-2 py-1 hover:bg-gray-100 cursor-pointer"
            >
              <input
                type="checkbox"
                checked={selected.includes(option.id)}
                onChange={() => toggleSelect(option.id)}
                className="mr-2"
              />
              {option.name}
            </label>
          ))}
        </div>
      )}
    </div>
  );
}

export default MultiSelectDropdown;
