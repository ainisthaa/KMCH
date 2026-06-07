function PriceRangeSelect({ setMinPrice, setMaxPrice, setPage }) {
    return (
      <div>
        <p>Price Range</p>
        <select
          className="bg-amber-200"
          onChange={(e) => {
            const [min, max] = e.target.value.split("-");
            setMinPrice(min);
            setMaxPrice(max);
            setPage(1);
          }}
        >
          <option value="">-- ทั้งหมด --</option>
          <option value="0-10">0 - 10</option>
          <option value="10-20">10 - 20</option>
          <option value="20-30">20 - 30</option>
          <option value="30-40">30 - 40</option>
          <option value="40-50">40 - 50</option>
          <option value="50-1000">50+</option>
        </select>
      </div>
    );
  }
  
  export default PriceRangeSelect;
  