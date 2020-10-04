#! /bin/bash

julia -p auto -t auto -q --project=. -e 'using PackageCompiler;\
    create_app("", "app", app_name="geoclose", force=true, filter_stdlibs=true, precompile_execution_file="precompile_app.jl")'
