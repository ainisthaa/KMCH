import React, { useEffect, useState } from "react";
import Select from "react-select";


export default function IngredientMultiSelect({ value, onChange, onResetPage }) {
  const [options, setOptions] = useState([]);

  /* โหลดลิสต์วัตถุดิบครั้งเดียว */
  useEffect(() => {
    fetch("http://localhost:8890/ingredient?limit=1000")
      .then((res) => res.json())
      .then((res) =>
        setOptions(
          res.data.map((i) => ({
            value: i.name.toLowerCase(),
            label: i.name,
          }))
        )
      );
  }, []);

  return (
    <Select
      options={options}  
      isMulti
      placeholder="เลือกวัตถุดิบ"
      value={value}       
      onChange={(vals) => {
        onChange(vals);          
        onResetPage?.();       
      }}
    />
  );
}
