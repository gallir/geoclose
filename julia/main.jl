module Geoclose

using ArgParse
using DataFrames
using CSV
using Distances
import Base.Threads.@spawn

const stopAtDistance = 0.001

function main()
    args = parse_commandline()
    data = load_csv(args["data"])
    toSearch = load_csv(args["search"])
    res = processParallel(data, toSearch)
    if isempty(args["output"])
        print(res)
    else
        CSV.write(args["output"], res)
    end
end

function processParallel(data, toSearch)
    rows = nothing
    p = Threads.nthreads()
    println("Threads $p")
    s = size(toSearch, 1) ÷ p
    start = 1

    segments = []
    dispatched = []
    # Process segments en parallel
    while (start <= size(toSearch, 1))
        e = min(start + s - 1, size(toSearch, 1))
        v = @view toSearch[start:e, :]
        push!(segments, @view toSearch[start:e, :])
        println("start: $start, end: $e size: ", size(v, 1))
        start += s
    end

    for s in segments
        # Dispatch a segment
        t = @spawn process(data, s)
        push!(dispatched, t)
    end

    println("Dispatched ", length(dispatched))
    flush(stdout)

    # Collect the results
    for t in dispatched
        r = fetch(t)
        if isnothing(rows)
            rows = r
        else
            append!(rows, r)
        end
    end
    return rows
end

function process(data, toSearch)
    df = DataFrame(id_search = Int[], id_data = Int[], distance_m = Int[])
    for s in eachrow(toSearch)
        minDistance = maxintfloat()
        picked = nothing
        for d in eachrow(data)
            if d.latitude == 0.0 && d.longitude == 0.0
                continue
            end
            if abs(d.latitude - s.latitude) > 1 || abs(d.longitude - s.longitude) > 1
                continue
            end
            dist = haversine((d.latitude, d.longitude), (s.latitude, s.longitude), 6371)
            if dist < minDistance
                minDistance = dist
                picked = d
                if dist < stopAtDistance
                    break
                end
            end
        end
        if isnothing(picked)
            continue
        end
        t = (s.id, picked.id, Int(floor(Int, minDistance * 1000)))
        # println(t)
        push!(df, t)
    end
    return df
end

function load_csv(filename)
    df = CSV.read(filename, DataFrame)
    return df
end

function parse_commandline()
    s = ArgParseSettings()

    @add_arg_table! s begin
        "--data", "-d"
        help = "Data file, for example giata.csv"
        arg_type = String
        required = true
        "--search", "-s"
        help = "Data to look for example new_properties.csv"
        arg_type = String
        required = true
        "--output", "-o"
        help = "Output CSV filename, if not specified, it prints in stdout"
        arg_type = String
        required = false
    end

    return parse_args(s)
end

main()
end