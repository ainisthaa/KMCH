import React from 'react';

function Pagination({ page, setPage, rowsPerPage, setRowsPerPage, total }) {
  const totalPages = Math.ceil(total / rowsPerPage);

  return (
    <div className="flex items-center gap-4 mt-4 flex-wrap justify-center">
      <button
        onClick={() => setPage(page - 1)}
        disabled={page === 1}
        className="px-3 py-1 bg-gray-200 rounded disabled:opacity-50"
      >
        &lt; Prev
      </button>

      <span>Page {page} / {totalPages}</span>

      <button
        onClick={() => setPage(page + 1)}
        disabled={page === totalPages}
        className="px-3 py-1 bg-gray-200 rounded disabled:opacity-50"
      >
        Next &gt;
      </button>

      <label className="font-bold ml-4">แสดงต่อหน้า:</label>
      <select
        className="border rounded px-2 py-1"
        value={rowsPerPage}
        onChange={(e) => {
          setRowsPerPage(Number(e.target.value));
          setPage(1);
        }}
      >
        <option value={5}>5</option>
        <option value={10}>10</option>
        <option value={15}>15</option>
      </select>
    </div>
  );
}

export default Pagination;
