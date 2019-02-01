# Domain Security

## Opis problemu
W dużych przedsiębiorstwach zarządzanie dostępem do witryn jest podstawowym narzędziem poprawiającym cyberbezpieczeństwo. Wraz z rozrostem takich list coraz ciężej jest zachować ich aktualność. Cyberprzestępcy mogą to wykorzystać, wykupując domeny o niegdyś dobrej reputacji. Niezbędy jest system, który będzie w stanie monitorować informacje o bezpieczeńswie stron, najlpiej agregująć wiele źródeł.

## Opis rozwiązania
Aplikacja służy do zarządzania whitelistami i blacklistami z domenami. Globalna lista śledzi wszystkie zapisane domeny i na bierząco aktualizuje informacje o ich bezpieczeństwie. Robi to na podstawie informacji o lukach w zabezpieczeniach z serwisów Talos i Shodan. Dzieki temu użytkownicy aplikacji mają bierzący podgląd na bezpieczeństwo stron blokowanych i przepuszczanych. Dzięki możliwości dowolnego kształtowania formatu wyjściowego list, mogą być one przystosowane zarówno do przetworzenia przez systemy komputerowe, jak i człowieka (Np. w formie pliku html z linkami). System umożlwia mapowanie wiele do wielu, co pozwala zarówno serwować jedną listę w wielu formatach, jak i wiele list w jednym formacie. Umożliwia o minimalizacje redundancji.

## Budowanie
Aby aplikacja dzialala, niezbedne jest srodowisko python 3 z blibliotekami:
- pymysql
- talos (Biblioteka ulatwiajaca parsowanie danych z talosinteligence) 
- shodan (analogiczna dla serwiso shodan.io).
Projekt jest napisany w jezykach go (backend), mysql (baza danych) ,html5,css,js (frontend) i python (parsowanie stron z informacjami o niebezpieczenstwie). 
Aby skompilwoac projekt nalezy wywolac komende 
```bash
[admin@linux threataway]$ go build
```
z katalogu glownego projektu. Należy mieć na uwadze, żeby GO_PATH kierowało na ten właśnie folder. 

## Uruchamianie
Serwer uruchamia się ze zbudowanego pliku wykonywalnego threataway
```bash
[admin@linux threataway]$ ./threataway
```
