using Geoclose
# include("src/Geoclose.jl")

path = dirname(Base.source_path())
append!(ARGS, ["-s", "$path/csv/search.csv", "-d", "$path/csv/data.csv", "-o", "/tmp/geoclose_precompile.csv"]) 
Geoclose.julia_main()
