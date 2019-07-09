# envsub

Substitute Environment variable definition in template file by real values of ENV

## Examples of template variable definition:
```
 ${VARIABLE}            Value of $VARIABLE or empty if variable is not set.
 ${VARIABLE=default}    Value of $VARIABLE or "default" if variable is not set.
 ${VARIABLE-}           Value of $VARIABLE or skip whole line if variable is not set.
 ```

## Usage

Print to console
```
$ envsub -i example.tpl
```

Save to file
```
$ envsub -i example.tpl > output.conf
```
