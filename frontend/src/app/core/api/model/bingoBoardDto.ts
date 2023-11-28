import {BingoCardDto} from './bingoCardDto';

export class BingoBoardDto {
    'name': string;
    'short_uuid': string;
    'bingo_cards': BingoCardDto[];
}
