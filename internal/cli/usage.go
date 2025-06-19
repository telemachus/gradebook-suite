package cli

var (
	calcUsage = `usage: gradebook-calc [-class CLASS -dir DIR -term TERM] [-help -version]

Calculate and print the grades for a class

options:
    -class CLASS  Class file to use (default: ./class.json)
    -dir DIR      Directory for gradebook and class.json files (default: ".")
    -term TERM    Limit calculation to grades in a given TERM

general:
    -help         Print this message
    -version      Print version`

	emailsUsage = `usage: gradebook-emails [-class CLASS -dir DIR] [-help -version]

Print the emails of students in a class

options:
    -class CLASS  Class file to use (default: ./class.json)
    -dir DIR      Directory for gradebook and class.json files (default: ".")

general:
    -help         Print this message
    -version      Print version`

	namesUsage = `usage: gradebook-names [-class CLASS -dir DIR -last-first] [-help -version]

Print the names of students in a class (in "First Last" or "Last, First" format)

options:
    -class CLASS  Class file to use (default: $PWD/class.json)
    -dir DIR      Directory for gradebook and class.json files (default: $PWD)
    -last-first   Print names in "Last, First" format (default: "First Last")

general:
    -help         Print this message
    -version      Print version`

	newUsage = `usage: gradebook-new -name NAME -type TYPE [-class CLASS -date DATE -dir DIR] [-help -version]

Create a new gradebook file for a class

required flags:
    -name NAME    Name of the gradebook file (only [A-Za-z0-9._-] are valid)
    -type TYPE    Type of gradebook file to create (must be in class.json)

options:
    -class        Class file to use (default: ./class.json)
    -date DATE    YYYYMMDD date for gradebook file (default: current date)
    -dir DIR      Directory for gradebook and class.json files (default: $PWD)

general:
    -help         Print this message
    -version      Print version`
)
