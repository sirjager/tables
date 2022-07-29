_red='\033[031m'
_grn='\033[032m'
_blu='\033[034m'
_nc='\033[0m' # No Color
number_regex='^[0-9]+$'


script_path="migrate"
migration_files_dir="./migrations"
database_url="postgresql://postgres:password@localhost:5432/gobase?sslmode=disable"

total_migration_files="$(find $migration_files_dir -type f | wc -l)"
total_migrations=$(($total_migration_files / 2)) 


printf "\n${_grn}Database Migration cli${_nc}"
printf "\n${_grn}1. ${_blu}Upgrade database to given level${_nc}"
printf "\n${_grn}2. ${_blu}Downgrade database to given level${_nc}"
printf "\n${_grn}3. ${_blu}Upgrade database to latest version${_nc}"
printf "\n${_grn}4. ${_red}Downgrade database to first version${_nc}. (-force) Skip Confirmation and downgrade."
printf "\n${_grn}5. ${_blu}Create new migration${_nc}"
printf "\n${_grn}6. ${_blu}List all migrations${_nc}"
printf "\n${_grn}7. ${_blu}Delete existing migration${_nc}"
printf "\n${_grn}8. ${_blu}Show current migration version${_nc}"
printf "\n${_grn}9. ${_blu}Set version V but don't run migration (ignores dirty state)${_nc}"
printf "\n${_grn}10. ${_red}Drop everything inside database. (-force) Skip Confirmation and delete everything.${_nc}"

if [[ "$1" == "" ]]; then
    printf "\n\n${_blu}Enter Option:${_nc} : ${_grn}"
    read -r option
else
    option=$1
    printf "\n\n${_blu}"
fi

if [[ "$option" == "3" ]]; then
    $script_path -database $database_url -path $migration_files_dir up
    printf "${_nc}"
    exit 0
elif [[ "$option" == "1" ]]; then
    if [[ "$2" != "" ]]; then
        level=$2
    else
        printf "\n${_blu}How many level you want to upgrade : ${_nc}"
        read -r level
    fi
    if ! [[ $level =~ $number_regex ]];  then
        printf "\n${_blu} $level ${_red}is not a valid number.$_nc}"
        printf "${_nc}"
        exit 1
    else
        printf "\n${_grn}Upgrading databaes to $level level \n"
        $script_path -database $database_url -path $migration_files_dir up $level
        printf "${_nc}"
        exit 0
    fi
elif [[ "$option" == "4" ]]; then
    if [[ "$2" == "-force" ]]; then
        $script_path -database $database_url -path $migration_files_dir down -all
    else
        $script_path -database $database_url -path $migration_files_dir down
    fi
    printf "${_nc}"
    exit 0
elif [[ "$option" == "2" ]]; then
    if [[ "$2" != "" ]]; then
        level=$2
    else
        printf "\n${_blu}How many level you want to downgrade : ${_nc}"
        read -r level
    fi
    if ! [[ $level =~ $number_regex ]]; then
        printf "\n${_blu} $level ${_red} is not a valid number : ${_nc}"
        printf "${_nc}"
        exit 1
    else
        printf "\n${_grn}Downgrading databaes to $level level \n"
        $script_path -database $database_url -path $migration_files_dir down $level
        printf "${_nc}"
        exit 0
    fi
elif [[ "$option" == "5" ]]; then
    if [[ "$2" != "" ]]; then
        title=$2
    else
        printf "\n${_blu}Migration title : ${_grn}"
        read -r title
    fi
    if [[ "$title" == "" ]] || [[ "$title" == " " ]]; then
        printf "\n${_red} No title specified. Exiting.\n"
        printf "${_nc}"
        exit 1
    else
        $script_path create -ext sql -dir $migration_files_dir -seq $title
        printf "${_nc}"
        exit 0
    fi
elif [[ "$option" == "6" ]]; then
    lsd --tree $migration_files_dir
    printf "\n${_blu}Total Migrations files :${_grn} ${total_migration_files}${_nc}"
    printf "\n${_blu}Total Migrations :${_grn} ${total_migrations}${_nc}\n"
    printf "${_nc}"
    exit 0
elif [[ "$option" == "8" ]]; then
    $script_path -source file:$migration_files_dir -database "$database_url" version
    printf "${_nc}"
    exit 0
elif [[ "$option" == "9" ]] || [[ "$1" == "setversion" ]];  then
    if [[ "$2" != "" ]]; then
        set_to_version=$2
    else
        printf "\n${_blu}Enter the version you want to force set : ${_grn}"
        read -r set_to_version
    fi
    if ! [[ $set_to_version =~ $number_regex ]]; then
        printf "\n${_blu} $version ${_red} is not a valid version number : ${_nc}"
        printf "${_nc}"
        exit 1
    else
        $script_path -source file:$migration_files_dir -database "$database_url" force $set_to_version
        printf "${_nc}"
        exit 0
    fi
elif [[ "$option" == "10" ]]; then
    if [[ "$2" == "-force" ]]; then
        $script_path -source file:$migration_files_dir -database "$database_url" drop -f
        printf "${_grn} Dropped everything inside the database.${_nc}\n\n"
    else
        $script_path -source file:$migration_files_dir -database "$database_url" drop
    fi
    printf "${_nc}"
    exit 0
elif [[ "$option" == "7" ]]; then
    printf "\n${_red} Not Implemented yet. Delete manually form $migration_files_dir\n"
    printf "${_nc}"
    exit 1
else
    printf "\n${_red} Not a valid option. Exiting.\n"
    printf "${_nc}"
    exit 1
fi

