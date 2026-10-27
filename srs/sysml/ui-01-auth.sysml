package UI_01_Auth {
    private import UiKit::*;
    part <auth> 'Вход и регистрация' : Row {
        doc /* REQ-FUNC-001, REQ-FUNC-002 */
        part 'Вход' : Screen {
            part 'Email [student@example.by]' : Input;
            part 'Пароль [скрыт]' : Input;
            part 'Войти' : Button;
            part 'Создать учётную запись' : Link;
            part 'Неверный email или пароль' : Banner;
            part 'После 5 ошибок — блокировка на 15 минут' : Label;
        }
        part 'Регистрация' : Screen {
            part 'Email [student@example.by]' : Input;
            part 'Пароль [12–128 символов]' : Input;
            part 'Зарегистрироваться' : Button;
            part 'Уже есть учётная запись? Войти' : Link;
            part 'Адрес занят или введён некорректно' : Banner;
            part 'После регистрации — настройка бюджета' : Label;
        }
    }
}
