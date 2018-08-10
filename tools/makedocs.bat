@echo off
SET "repo=.."
cd %repo%
pwd

SET "package=api"
echo Generating README.md for %package%
godocdown .\%package%\ > .\%package%\README.md

SET "package=application"
echo Generating README.md for %package%
godocdown .\%package%\ > .\%package%\README.md

SET "package=config"
echo Generating README.md for %package%
godocdown .\%package%\ > .\%package%\README.md

REM This wont work like the others, cant find package. So have to manually switch into and out of database folder.
REM godocdown Bug? or Microsoft Powershell? WTF!
SET "package=database"
cd %package%
echo Generating README.md for %package%
godocdown > README.md
cd ..

SET "package=models"
echo Generating README.md for %package%
godocdown .\%package%\ > .\%package%\README.md

SET "package=site"
echo Generating README.md for %package%
godocdown .\%package%\ > .\%package%\README.md

echo Done.