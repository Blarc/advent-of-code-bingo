import {BingoCardDto} from './bingoCardDto.model';

export class BingoBoardDto {
    'name': string;
    'short_uuid': string;
    'bingo_cards': BingoCardDto[];
}
