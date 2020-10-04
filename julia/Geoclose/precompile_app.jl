# using Geoclose
include("src/Geoclose.jl")

append!(ARGS, ["-s", "../../csv/search.csv", "-d", "../../csv/data.csv", "-o", "/tmp/geoclose_preconpile.csv"]) 
Geoclose.julia_main()