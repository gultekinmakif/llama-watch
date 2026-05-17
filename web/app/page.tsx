import { loadSnapshot } from '../lib/snapshot'

export default function Page() {
  const snapshot = loadSnapshot()
  const { columns, rows } = snapshot

  return (
    <main className="p-6">
      <h1 className="text-xl font-semibold mb-4">llama-watch</h1>
      <table className="border-collapse text-sm">
        <thead>
          <tr>
            <th className="border px-2 py-1 text-left">slug</th>
            <th className="border px-2 py-1 text-left">name</th>
            <th className="border px-2 py-1 text-left">category</th>
            <th className="border px-2 py-1 text-left">chains</th>
            {columns.map((column) => (
              <th key={column.key} className="border px-2 py-1 text-left">
                {column.label}
              </th>
            ))}
          </tr>
        </thead>
        <tbody>
          {rows.map((row) => (
            <tr key={row.slug}>
              <td className="border px-2 py-1">{row.slug}</td>
              <td className="border px-2 py-1">{row.name}</td>
              <td className="border px-2 py-1">{row.category ?? ''}</td>
              <td className="border px-2 py-1">{row.chains.join(', ')}</td>
              {columns.map((column) => (
                <td key={column.key} className="border px-2 py-1">
                  {String(row.cells[column.key])}
                </td>
              ))}
            </tr>
          ))}
        </tbody>
      </table>
    </main>
  )
}
