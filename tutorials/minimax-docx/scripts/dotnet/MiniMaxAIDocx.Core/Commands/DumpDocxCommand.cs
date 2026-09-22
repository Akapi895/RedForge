using System;
using System.IO;
using System.Linq;
using DocumentFormat.OpenXml.Packaging;
using DocumentFormat.OpenXml.Wordprocessing;

namespace MiniMaxAIDocx.Core.Commands;

public static class DumpDocxContent
{
    public static void Dump(string inputPath)
    {
        using var doc = WordprocessingDocument.Open(inputPath, false);
        var body = doc.MainDocumentPart.Document.Body;

        Console.WriteLine("=== DOCUMENT HEADINGS & SECTIONS ===");
        foreach (var p in body.Elements<Paragraph>())
        {
            var styleId = p.ParagraphProperties?.ParagraphStyleId?.Val?.Value;
            if (styleId != null && (styleId.StartsWith("Heading") || styleId == "Title" || styleId == "Subtitle"))
            {
                Console.WriteLine($"{styleId}: {p.InnerText}");
            }
        }

        Console.WriteLine("\n=== TABLES SUMMARY ===");
        int tCount = 0;
        foreach (var tbl in body.Elements<Table>())
        {
            tCount++;
            var rows = tbl.Elements<TableRow>().ToList();
            Console.WriteLine($"Table #{tCount}: {rows.Count} rows");
            for (int r = 0; r < Math.Min(3, rows.Count); r++)
            {
                var cells = rows[r].Elements<TableCell>().Select(c => c.InnerText.Trim()).ToList();
                Console.WriteLine($"  Row {r}: {string.Join(" | ", cells)}");
            }
        }
    }
}
