import {Pipe, PipeTransform} from '@angular/core';

import {BingoCardDto} from '../api/model/bingoCardDto';

@Pipe({
    standalone: true,
    name: 'userCount',
    pure: true
})
export class UserCountPipe implements PipeTransform {
    transform(card: BingoCardDto): string {
        if (card.selected) {
            const numOfOthers = card.user_count - 1;
            switch (numOfOthers) {
                case 0:
                    return `You and none other`;
                case 1:
                    return `You and one other`;
                default:
                    return `You and ${numOfOthers} others`;
            }
        } else {
            switch (card.user_count) {
                case 0:
                    return `No one`;
                case 1:
                    return `Only one other`;
                default:
                    return `${card.user_count} others`;
            }
        }
    }
}
