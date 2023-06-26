package main

func partial[A any, B any, R any](f func(A,B) R, a A) func(B) R {
  return func(b B) R { return f(a,b) }
}

