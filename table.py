import colored

def draw_table(headers=None, rows=None):
    # Default headers and rows if none are provided
    if headers is None:
        headers = ["Name", "Age", "Occupation"]
    if rows is None:
        rows = [
            ["Alice", "30", "Engineer"],
            ["Bob", "25", "Designer"],
            ["Charlie", "35", "Teacher"]
        ]
    
    # Determine the width for each column
    col_widths = [max(len(str(item)) for item in column) for column in zip(*([headers] + rows))]
    
    # ANSI escape for styling
    header_style = colored.fg("cyan") + colored.attr("bold")
    row_style = colored.fg("yellow")
    reset_style = colored.attr("reset")

    # Print header
    header_line = " | ".join(header_style + str(header).ljust(width) for header, width in zip(headers, col_widths)) + reset_style
    print(header_line)
    print("-" * (sum(col_widths) + 3 * (len(headers) - 1)))  # Divider line

    # Print rows
    for row in rows:
        row_line = " | ".join(row_style + str(cell).ljust(width) for cell, width in zip(row, col_widths)) + reset_style
        print(row_line)

# Call the function without arguments
draw_table()
