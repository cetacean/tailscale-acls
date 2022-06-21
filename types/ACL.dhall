let Action = ./Action.dhall

in  { Type = { action : Action, src : List Text, dst : List Text }
    , default =
      { action = Action.accept, src = [] : List Text, dst = [] : List Text }
    }
