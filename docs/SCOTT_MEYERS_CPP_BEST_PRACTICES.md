# Scott Meyers C++ Best Practices

## Sources
- **[EC]** Effective C++, 3rd Edition (55 Items)
- **[MEC]** More Effective C++ (35 Items)
- **[EMC]** Effective Modern C++: C++11 and C++14 (42 Items)

---

## Tier 1 — Safety & Correctness (Apply First)

These practices prevent crashes, memory leaks, undefined behavior, and subtle bugs. Adopt them before any performance tuning.

### Resource Management
- **Manage resources with RAII** — Wrap raw resources (memory, file handles, sockets, locks) in objects whose destructors release them. This prevents leaks, especially when exceptions are thrown. [EC Item 13]
- **Prevent resource leaks in constructors** — Initialize members in the constructor initialization list and use RAII members so that if a constructor throws, already-constructed members are cleaned up automatically. [EC Item 13]
- **Understand special member functions** — Know when the compiler generates default, copy, move, and destructor functions, and when you must suppress or define them yourself. [EMC Item 17]

### Polymorphism & Inheritance
- **Declare destructors virtual in polymorphic base classes** — Ensures derived-class destructors run when deleting through a base pointer. Without it, behavior is undefined. [EC Item 7]
- **Never call virtual functions during construction or destruction** — The virtual table is not yet fully set up (or is being torn down), so calls resolve to the current class's version, not the derived class's. [EC Item 9]
- **Distinguish inheritance of interface from inheritance of implementation** — Public inheritance means "is-a" (interface inheritance). If you only want implementation reuse, prefer composition or private inheritance. [MEC Item 33]
- **Consider alternatives to virtual functions** — Prefer function objects, `std::function`, or the Strategy pattern when runtime polymorphism is not strictly required. [MEC Item 35]

### Initialization & Copying
- **Initialize all members in the constructor initialization list** — More efficient than assignment in the body; required for `const`, reference, and non-default-constructible members. [EC Item 4]
- **Declare copying functions if compiler-generated ones are wrong** — If your class manages resources, explicitly define or delete the copy constructor, copy assignment, move constructor, and move assignment (Rule of Five). [EC Item 11]
- **Think carefully about copying behavior** — Decide whether your type should be copyable, movable, both, or neither, and express that decision explicitly. [EC Item 14]
- **Use `explicit` on single-argument constructors** — Prevents unintended implicit conversions that can lead to surprising behavior. [EC Item 22]

### Exception Safety
- **Write exception-safe code** — Provide at least the basic guarantee (no leaks, object remains valid), preferably the strong guarantee (commit-or-rollback), or no-throw where possible. [EC Items 29–32]
- **Treat exceptions as part of the interface** — Document which functions can throw and what exceptions they may throw, so callers can write correct `try`/`catch` logic. [EC Item 29]

### Const Correctness
- **Use `const` whenever possible** — `const` member functions, `const` references, and `const` variables communicate intent, enable compiler optimizations, and prevent accidental mutation. [EC Item 3]
- **Prefer `const_iterators` to iterators** — When you only need to read, use `const_iterator` to communicate read-only intent and avoid accidental modification. [EMC Item 13]

### Modern C++ Safety
- **Prefer `nullptr` to `0` and `NULL`** — Avoids ambiguity with integer overloads and is type-safe. [EMC Item 8]
- **Prefer scoped enums to unscoped enums** — `enum class` prevents name pollution and implicit conversion to integers. [EMC Item 25]
- **Prefer `auto` to explicit type declarations** — Reduces verbosity, avoids type mismatches, and works well with templates and lambdas. [EMC Item 26]
- **Declare overridden functions `override`** — Makes intent explicit and lets the compiler catch signature mismatches. [EMC Item 11]
- **Declare functions `noexcept` if they won't throw** — Enables move optimizations and documents exception guarantees. [EMC Item 14]

---

## Tier 2 — Performance & STL Usage

These practices improve day-to-day performance, reduce allocations, and make standard library usage idiomatic.

### Containers & Memory
- **Prefer `std::vector` and `std::string` by default** — Contiguous storage is cache-friendly, fast to traverse, and usually the right default unless you have a proven need for `list`, `deque`, `set`, etc. [EMC Items 23–24]
- **Prefer `std::array` over built-in arrays** — Fixed-size stack container that knows its size and does not decay to a pointer. [EMC Item 23]
- **Use `reserve()` to avoid reallocations** — If you know the approximate size of a `vector` or `string`, reserve capacity up front to avoid repeated allocations and copies. [EC Item 17 / EMC Item 42]
- **Use `std::unique_ptr` for exclusive ownership** — Expresses single-owner semantics and avoids manual `delete`. [EMC Item 18]
- **Use `std::shared_ptr` for shared ownership** — For cases where multiple owners genuinely need shared lifetime management. [EMC Item 19]
- **Use `std::weak_ptr` for `std::shared_ptr`-like pointers that can dangle** — Breaks reference cycles and observes shared objects without extending their lifetime. [EMC Item 20]
- **Prefer `std::make_unique` and `std::make_shared` to direct use of `new`** — Exception-safe, more concise, and (for `make_shared`) more efficient. [EMC Item 21]

### Algorithms & Iterators
- **Prefer algorithms to hand-written loops** — Standard algorithms are less error-prone, often optimized, and communicate intent more clearly than raw `for`/`while` loops. [EC Item 43]
- **Use the right container for the job** — `vector` for general use; `list` only when you need stable iterators and frequent middle insertions; `deque` for front/back insertion; `set`/`map` for sorted unique keys; `unordered_*` for hash-based lookups. [EC Items 16–23]
- **Prefer member functions over algorithms when faster** — Some containers provide specialized member functions (e.g., `list::remove`, `list::sort`, `map::find`) that outperform generic algorithm equivalents. [EC Item 44]
- **Be careful with iterator invalidation** — Insertions and erasures can invalidate iterators, pointers, and references. Re-obtain them after mutating a container. [EC Item 12]
- **Prefer `++it` over `it++` for iterators** — Pre-increment avoids creating a temporary object; for iterators this can be a measurable performance win. [EC Item 6]

### Move Semantics & Perfect Forwarding
- **Consider pass by value for copyable parameters that are cheap to move and always copied** — Enables a single function to handle both lvalues and rvalues efficiently. [EMC Item 41]
- **Consider emplacement instead of insertion** — `emplace_back`, `emplace`, etc., construct elements in-place and can avoid temporary objects. [EMC Item 42]
- **Use `std::forward` when perfect forwarding** — Preserves the value category of forwarded arguments in generic wrappers. [EMC Item 40]

---

## Tier 3 — Design & Maintainability

These practices refine API design, reduce coupling, and improve long-term maintainability.

### Interface Design
- **Make interfaces easy to use correctly and hard to use incorrectly** — Good API design prevents misuse at compile time (e.g., strong types, `explicit`, deleted functions) rather than relying on documentation. [EC Item 18]
- **Prefer pass-by-reference-to-const to pass-by-value** — Avoids unnecessary copies and prevents slicing for user-defined types. [EC Item 20]
- **Make data members `private`** — Preserves encapsulation and class invariants; non-member functions can often provide the same functionality with better encapsulation. [EC Item 22]
- **Prefer non-member, non-friend functions** — Functions that do not need access to private members should be non-member, non-friend functions to increase encapsulation and reduce coupling. [EC Item 23]

### Composition vs. Inheritance
- **Avoid unnecessary inheritance** — Prefer composition ("has-a") over inheritance ("is-a") unless you truly need polymorphism. Inheritance couples classes tightly. [MEC Item 33]
- **Prefer composition to inheritance** — When modeling "has-a" or "is-implemented-in-terms-of" relationships, composition yields looser coupling and easier testing. [MEC Item 33]

### Modern C++ Style
- **Prefer alias declarations to `typedefs`** — `using` is more readable, works with templates, and supports alias templates. [EMC Item 9]
- **Distinguish between `()` and `{}` when creating objects** — Brace initialization prevents narrowing conversions and is more consistent across types. [EMC Item 7]
- **Use the Pimpl Idiom judiciously** — Hides implementation details in the `.cpp` file to reduce compile times and preserve binary compatibility. [EMC Item 22]

---

## Practical Adoption Order

If you are introducing these practices to a team or codebase, apply them in this order:

1. **Tier 1 first** — Safety and correctness. These prevent crashes, memory leaks, and undefined behavior.
2. **Tier 2 next** — Performance and STL usage. These improve day-to-day efficiency and standard library idioms.
3. **Tier 3 last** — Design and maintainability. These refine APIs and long-term code health.

Start with the highest-impact, lowest-risk items: RAII, virtual destructors, `const`, `nullptr`, `override`, `std::vector`, and `std::unique_ptr`.
