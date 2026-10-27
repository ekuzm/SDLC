// Общие элементы схематических экранов «Баланс».
package UiKit {
    abstract part def Widget;
    part def Screen :> Widget;
    part def Column :> Widget;
    part def Row :> Widget;
    part def Form :> Column;
    part def ListView :> Column;
    part def ListItem :> Row;
    part def Label :> Widget;
    part def Input :> Widget;
    part def Select :> Widget;
    part def Button :> Widget;
    part def Link :> Widget;
    part def Chart :> Widget;
    part def Banner :> Widget;
    part def NavBar :> Row;
}
