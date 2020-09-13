compdef _gcmd gcmd

function _gcmd {
    _arguments "1: :($(gcmd -ls))"
}
