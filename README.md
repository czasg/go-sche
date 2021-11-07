# go-sche

```text
|—————————————|                          notify
|  scheduler  | ————————> task<labels> |-------> label |-----> handler
|—————————————|                        |               |-----> handler
       |  interface                    |               |-----> handler
|—————————————|                        |
|    store    |                        |-------> label |-----> handler
|—————————————|                        |               | ...
       |  next run time                | ...
|—————————————|
|    task     |
|—————————————|
       |  last run time
|—————————————|
|  cron-trig  |
|—————————————|
```